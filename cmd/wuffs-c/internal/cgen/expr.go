// Copyright 2017 The Wuffs Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cgen

import (
	"fmt"
	"strings"

	"github.com/google/wuffs/lang/builtin"

	a "github.com/google/wuffs/lang/ast"
	t "github.com/google/wuffs/lang/token"
)

func (g *gen) writeExpr(b *buffer, n *a.Expr, rp replacementPolicy, pp parenthesesPolicy, depth uint32) error {
	if depth > a.MaxExprDepth {
		return fmt.Errorf("expression recursion depth too large")
	}
	depth++

	if rp == replaceCallSuspendibles && n.CallSuspendible() {
		if g.currFunk.tempR >= g.currFunk.tempW {
			return fmt.Errorf("internal error: temporary variable count out of sync")
		}
		// TODO: check that this works with nested call-suspendibles:
		// "foo?().bar().qux?()(p?(), q?())".
		//
		// Also be aware of evaluation order in the presence of side effects:
		// in "foo(a?(), b!(), c?())", b should be called between a and c.
		b.printf("%s%d", tPrefix, g.currFunk.tempR)
		g.currFunk.tempR++
		return nil
	}

	if cv := n.ConstValue(); cv != nil {
		if !n.MType().IsBool() {
			b.writes(cv.String())
		} else if cv.Cmp(zero) == 0 {
			b.writes("false")
		} else if cv.Cmp(one) == 0 {
			b.writes("true")
		} else {
			return fmt.Errorf("%v has type bool but constant value %v is neither 0 or 1", n.Str(g.tm), cv)
		}
		return nil
	}

	switch op := n.Operator(); {
	case op.IsXUnaryOp():
		return g.writeExprUnaryOp(b, n, rp, pp, depth)
	case op.IsXBinaryOp():
		return g.writeExprBinaryOp(b, n, rp, pp, depth)
	case op.IsXAssociativeOp():
		return g.writeExprAssociativeOp(b, n, rp, pp, depth)
	}
	return g.writeExprOther(b, n, rp, pp, depth)
}

func (g *gen) writeExprOther(b *buffer, n *a.Expr, rp replacementPolicy, pp parenthesesPolicy, depth uint32) error {
	switch n.Operator() {
	case 0:
		if id1 := n.Ident(); id1 == t.IDThis {
			b.writes("self")
		} else {
			if n.GlobalIdent() {
				b.writes(g.pkgPrefix)
			} else {
				b.writes(vPrefix)
			}
			b.writes(id1.Str(g.tm))
		}
		return nil

	case t.IDOpenParen:
		// n is a function call.
		// TODO: delete this hack that only matches "foo.bar_bits(etc)".
		if isThatMethod(g.tm, n, t.IDLowBits, 1) {
			// "x.low_bits(n:etc)" in C is "((x) & ((1 << (n)) - 1))".
			x := n.LHS().Expr().LHS().Expr()
			b.writes("((")
			if err := g.writeExpr(b, x, rp, parenthesesOptional, depth); err != nil {
				return err
			}
			b.writes(") & ((1 << (")
			if err := g.writeExpr(b, n.Args()[0].Arg().Value(), rp, parenthesesOptional, depth); err != nil {
				return err
			}
			b.writes(")) - 1))")
			return nil
		}
		if isThatMethod(g.tm, n, t.IDHighBits, 1) {
			// "x.high_bits(n:etc)" in C is "((x) >> (8*sizeof(x) - (n)))".
			x := n.LHS().Expr().LHS().Expr()
			b.writes("((")
			if err := g.writeExpr(b, x, rp, parenthesesOptional, depth); err != nil {
				return err
			}
			b.writes(") >> (")
			if sz, err := g.sizeof(x.MType()); err != nil {
				return err
			} else {
				b.printf("%d", 8*sz)
			}
			b.writes(" - (")
			if err := g.writeExpr(b, n.Args()[0].Arg().Value(), rp, parenthesesOptional, depth); err != nil {
				return err
			}
			b.writes(")))")
			return nil
		}
		if isThatMethod(g.tm, n, t.IDIsError, 0) || isThatMethod(g.tm, n, t.IDIsOK, 0) ||
			isThatMethod(g.tm, n, t.IDIsSuspension, 0) {
			if pp == parenthesesMandatory {
				b.writeb('(')
			}
			x := n.LHS().Expr().LHS().Expr()
			if err := g.writeExpr(b, x, rp, parenthesesMandatory, depth); err != nil {
				return err
			}
			switch key := n.LHS().Expr().Ident(); key {
			case t.IDIsError:
				b.writes(" < 0")
			case t.IDIsOK:
				b.writes(" == 0")
			case t.IDIsSuspension:
				b.writes(" > 0")
			default:
				return fmt.Errorf("unrecognized token (0x%X) for writeExprOther's IsXxx", key)
			}
			if pp == parenthesesMandatory {
				b.writeb(')')
			}
			return nil
		}
		if isThatMethod(g.tm, n, t.IDSuffix, 1) {
			// TODO: don't assume that the slice is a slice of base.u8.
			b.writes("wuffs_base__slice_u8_suffix(")
			x := n.LHS().Expr().LHS().Expr()
			if err := g.writeExpr(b, x, rp, parenthesesOptional, depth); err != nil {
				return err
			}
			b.writeb(',')
			if err := g.writeExpr(b, n.Args()[0].Arg().Value(), rp, parenthesesOptional, depth); err != nil {
				return err
			}
			b.writes(")")
			return nil
		}
		if isInSrc(g.tm, n, t.IDLimit, 1) {
			return fmt.Errorf(`TODO: cgen an "in.src.limit" expression`)
		}
		if isInSrc(g.tm, n, t.IDMark, 0) {
			b.printf("wuffs_base__io_reader__mark(&%ssrc, %srptr_src)", aPrefix, bPrefix)
			return nil
		}
		if isInSrc(g.tm, n, t.IDSinceMark, 0) {
			b.printf("((wuffs_base__slice_u8){ "+
				".ptr = %ssrc.private_impl.mark, "+
				".len = %ssrc.private_impl.mark ? (size_t)(%srptr_src - %ssrc.private_impl.mark) : 0, })",
				aPrefix, aPrefix, bPrefix, aPrefix)
			return nil
		}
		// TODO: io_reader.is_marked, not just io_writer.is_marked?
		if isInDst(g.tm, n, t.IDLimit, 1) {
			return fmt.Errorf(`TODO: cgen an "in.dst.limit" expression`)
		}
		if isInDst(g.tm, n, t.IDMark, 0) {
			// TODO: is a private_impl.mark the right representation? What if
			// the function is passed a (ptr io_writer) instead of a
			// (io_writer)? Do we still want to have that mark live outside of
			// the function scope?
			b.printf("wuffs_base__io_writer__mark(&%sdst, %swptr_dst)", aPrefix, bPrefix)
			return nil
		}
		if isInDst(g.tm, n, t.IDSinceMark, 0) {
			// Write .len as either "foo ? bar : baz" or "bar".
			//
			// TODO: drop the "true" in the "if true", provided that the
			// benchmark numbers improve.
			len0, len1 := "", ""
			if true || !n.BoundsCheckOptimized() {
				len0 = aPrefix + "dst.private_impl.mark ?"
				len1 = ": 0"
			}
			b.printf("((wuffs_base__slice_u8){ "+
				".ptr = %sdst.private_impl.mark, "+
				".len = %s (size_t)(%swptr_dst - %sdst.private_impl.mark) %s, })",
				aPrefix, len0, bPrefix, aPrefix, len1)
			return nil
		}
		if isInDst(g.tm, n, t.IDIsMarked, 0) {
			if pp == parenthesesMandatory {
				b.writeb('(')
			}
			b.printf("%sdst.private_impl.mark != NULL", aPrefix)
			if pp == parenthesesMandatory {
				b.writeb(')')
			}
			return nil
		}
		if isInDst(g.tm, n, t.IDCopyFromReader32, 2) {
			b.printf("wuffs_base__io_writer__copy_from_reader32(&%swptr_dst, %swend_dst",
				bPrefix, bPrefix)
			// TODO: don't assume that the first argument is "in.src".
			b.printf(", &%srptr_src, %srend_src,", bPrefix, bPrefix)
			a := n.Args()[1].Arg().Value()
			if err := g.writeExpr(b, a, rp, parenthesesOptional, depth); err != nil {
				return err
			}
			b.writeb(')')
			return nil
		}
		if isInDst(g.tm, n, t.IDCopyFromHistory32, 2) {
			bco := ""
			if n.BoundsCheckOptimized() {
				bco = "__bco"
			}
			b.printf("wuffs_base__io_writer__copy_from_history32%s("+
				"&%swptr_dst, %sdst.private_impl.mark , %swend_dst",
				bco, bPrefix, aPrefix, bPrefix)
			for _, o := range n.Args() {
				b.writeb(',')
				if err := g.writeExpr(b, o.Arg().Value(), rp, parenthesesOptional, depth); err != nil {
					return err
				}
			}
			b.writeb(')')
			return nil
		}
		if isInDst(g.tm, n, t.IDCopyFromSlice32, 2) {
			b.printf("wuffs_base__io_writer__copy_from_slice32("+
				"&%swptr_dst, %swend_dst", bPrefix, bPrefix)
			for _, o := range n.Args() {
				b.writeb(',')
				if err := g.writeExpr(b, o.Arg().Value(), rp, parenthesesOptional, depth); err != nil {
					return err
				}
			}
			b.writeb(')')
			return nil
		}
		if isInDst(g.tm, n, t.IDCopyFromSlice, 1) {
			b.printf("wuffs_base__io_writer__copy_from_slice(&%swptr_dst, %swend_dst,", bPrefix, bPrefix)
			a := n.Args()[0].Arg().Value()
			if err := g.writeExpr(b, a, rp, parenthesesOptional, depth); err != nil {
				return err
			}
			b.writeb(')')
			return nil
		}
		if isThatMethod(g.tm, n, t.IDCopyFromSlice, 1) {
			b.writes("wuffs_base__slice_u8__copy_from_slice(")
			receiver := n.LHS().Expr().LHS().Expr()
			if err := g.writeExpr(b, receiver, rp, parenthesesOptional, depth); err != nil {
				return err
			}
			b.writeb(',')
			a := n.Args()[0].Arg().Value()
			if err := g.writeExpr(b, a, rp, parenthesesOptional, depth); err != nil {
				return err
			}
			b.writes(")\n")
			return nil
		}
		if isThatMethod(g.tm, n, t.IDLength, 0) {
			if pp == parenthesesMandatory {
				b.writeb('(')
			}
			b.writes("(uint64_t)(")
			if err := g.writeExpr(b, n.LHS().Expr().LHS().Expr(), rp, parenthesesMandatory, depth); err != nil {
				return err
			}
			b.writes(".len)")
			if pp == parenthesesMandatory {
				b.writeb(')')
			}
			return nil
		}
		if isThatMethod(g.tm, n, t.IDAvailable, 0) {
			if pp == parenthesesMandatory {
				b.writeb('(')
			}
			p0, p1 := "", ""
			if o := n.LHS().Expr().LHS().Expr(); o != nil {
				// TODO: don't hard-code these.
				switch o.Str(g.tm) {
				case "in.dst":
					p0 = bPrefix + "wend_dst"
					p1 = bPrefix + "wptr_dst"
				case "in.src":
					p0 = bPrefix + "rend_src"
					p1 = bPrefix + "rptr_src"
				}
			}
			if p0 == "" {
				return fmt.Errorf(`TODO: cgen a "foo.available" expression`)
			}
			b.printf("(uint64_t)(%s - %s)", p0, p1)
			if pp == parenthesesMandatory {
				b.writeb(')')
			}
			return nil
		}
		if isThatMethod(g.tm, n, g.tm.ByName("update"), 1) {
			// TODO: don't hard-code the class name or this.checksum.
			class := "wuffs_crc32__ieee_hasher"
			if g.pkgName == "zlib" {
				class = "wuffs_adler32__hasher"
			}
			b.printf("%s__update(&self->private_impl.f_checksum, ", class)
			a := n.Args()[0].Arg().Value()
			if err := g.writeExpr(b, a, rp, parenthesesMandatory, depth); err != nil {
				return err
			}
			b.writes(")\n")
			return nil
		}
		if isThatMethod(g.tm, n, g.tm.ByName("set_literal_width"), 1) {
			// TODO: don't hard-code lzw.
			b.printf("%slzw_decoder__set_literal_width(&self->private_impl.f_lzw, ", g.pkgPrefix)
			a := n.Args()[0].Arg().Value()
			if err := g.writeExpr(b, a, rp, parenthesesMandatory, depth); err != nil {
				return err
			}
			b.writes(")\n")
			return nil
		}
		if isThatMethod(g.tm, n, t.IDLimit, 1) {
			return fmt.Errorf(`TODO: cgen a "foo.limit" expression`)
		}
		if isThatMethod(g.tm, n, t.IDMark, 0) {
			// TODO: don't hard-code v_r or b_rptr_src.
			b.printf("wuffs_base__io_reader__mark(&v_r, b_rptr_src)")
			return nil
		}
		if isThatMethod(g.tm, n, t.IDSinceMark, 0) {
			// TODO: don't hard-code v_r or b_rptr_src.
			b.printf("((wuffs_base__slice_u8){ " +
				".ptr = v_r.private_impl.mark, " +
				".len = v_r.private_impl.mark ? (size_t)(b_rptr_src - v_r.private_impl.mark) : 0, })")
			return nil
		}
		if isThatMethod(g.tm, n, g.tm.ByName("initialize"), 5) {
			// TODO: don't hard-code a_dst.
			b.printf("wuffs_base__image_config__initialize(a_dst")
			for _, o := range n.Args() {
				b.writeb(',')
				if err := g.writeExpr(b, o.Arg().Value(), rp, parenthesesOptional, depth); err != nil {
					return err
				}
			}
			b.printf(")")
			return nil
		}
		// TODO.

	case t.IDOpenBracket:
		// n is an index.
		if err := g.writeExpr(b, n.LHS().Expr(), rp, parenthesesMandatory, depth); err != nil {
			return err
		}
		if lTyp := n.LHS().Expr().MType(); lTyp.IsSliceType() {
			// TODO: don't assume that the slice is a slice of base.u8.
			b.writes(".ptr")
		}
		b.writeb('[')
		if err := g.writeExpr(b, n.RHS().Expr(), rp, parenthesesOptional, depth); err != nil {
			return err
		}
		b.writeb(']')
		return nil

	case t.IDColon:
		// n is a slice.
		lhs := n.LHS().Expr()
		mhs := n.MHS().Expr()
		rhs := n.RHS().Expr()
		switch {
		case mhs != nil && rhs == nil:
			b.writes("wuffs_base__slice_u8__subslice_i(")
		case mhs == nil && rhs != nil:
			b.writes("wuffs_base__slice_u8__subslice_j(")
		case mhs != nil && rhs != nil:
			b.writes("wuffs_base__slice_u8__subslice_ij(")
		}

		lhsIsArray := lhs.MType().IsArrayType()
		if lhsIsArray {
			// TODO: don't assume that the slice is a slice of base.u8.
			b.writes("((wuffs_base__slice_u8){.ptr=")
		}
		if err := g.writeExpr(b, lhs, rp, parenthesesOptional, depth); err != nil {
			return err
		}
		if lhsIsArray {
			b.printf(",.len=%v})", lhs.MType().ArrayLength().ConstValue())
		}

		if mhs != nil {
			b.writeb(',')
			if err := g.writeExpr(b, mhs, rp, parenthesesOptional, depth); err != nil {
				return err
			}
		}
		if rhs != nil {
			b.writeb(',')
			if err := g.writeExpr(b, rhs, rp, parenthesesOptional, depth); err != nil {
				return err
			}
		}
		if mhs != nil || rhs != nil {
			b.writeb(')')
		}
		return nil

	case t.IDDot:
		lhs := n.LHS().Expr()
		if lhs.Ident() == t.IDIn {
			b.writes(aPrefix)
			b.writes(n.Ident().Str(g.tm))
			return nil
		}

		if err := g.writeExpr(b, lhs, rp, parenthesesMandatory, depth); err != nil {
			return err
		}
		if key := lhs.MType().Decorator(); key == t.IDPtr || key == t.IDNptr {
			b.writes("->")
		} else {
			b.writes(".")
		}
		b.writes("private_impl." + fPrefix)
		b.writes(n.Ident().Str(g.tm))
		return nil

	case t.IDError, t.IDStatus, t.IDSuspension:
		status := g.statusMap[n.StatusQID()]
		if status.name == "" {
			msg, _ := t.Unescape(n.Ident().Str(g.tm))
			z := builtin.StatusMap[msg]
			if z.Message == "" {
				return fmt.Errorf("no status code for %q", msg)
			}
			status.name = strings.ToUpper(g.cName(z.String()))
		}
		b.writes(status.name)
		return nil
	}
	return fmt.Errorf("unrecognized token (0x%X) for writeExprOther", n.Operator())
}

func (g *gen) writeExprUnaryOp(b *buffer, n *a.Expr, rp replacementPolicy, pp parenthesesPolicy, depth uint32) error {
	op := n.Operator()
	opName := cOpNames[0xFF&op]
	if opName == "" {
		return fmt.Errorf("unrecognized operator %q", op.AmbiguousForm().Str(g.tm))
	}

	b.writes(opName)
	return g.writeExpr(b, n.RHS().Expr(), rp, parenthesesMandatory, depth)
}

func (g *gen) writeExprBinaryOp(b *buffer, n *a.Expr, rp replacementPolicy, pp parenthesesPolicy, depth uint32) error {
	opName, tilde := "", false

	op := n.Operator()
	switch op {
	case t.IDXBinaryTildeSatPlus, t.IDXBinaryTildeSatMinus:
		uBits := uintBits(n.MType().QID())
		if uBits == 0 {
			return fmt.Errorf("unsupported tilde-operator type %q", n.MType().Str(g.tm))
		}
		uOp := "add"
		if op != t.IDXBinaryTildeSatPlus {
			uOp = "sub"
		}
		b.printf("wuffs_base__u%d__sat_%s", uBits, uOp)
		opName, tilde = ",", true

	case t.IDXBinaryAs:
		return g.writeExprAs(b, n.LHS().Expr(), n.RHS().TypeExpr(), rp, depth)

	default:
		opName = cOpNames[0xFF&op]
		if opName == "" {
			return fmt.Errorf("unrecognized operator %q", op.AmbiguousForm().Str(g.tm))
		}
	}

	if pp == parenthesesMandatory || tilde {
		b.writeb('(')
	}
	if err := g.writeExpr(b, n.LHS().Expr(), rp, parenthesesMandatory, depth); err != nil {
		return err
	}
	b.writes(opName)
	if err := g.writeExpr(b, n.RHS().Expr(), rp, parenthesesMandatory, depth); err != nil {
		return err
	}
	if pp == parenthesesMandatory || tilde {
		b.writeb(')')
	}
	return nil
}

func (g *gen) writeExprAs(b *buffer, lhs *a.Expr, rhs *a.TypeExpr, rp replacementPolicy, depth uint32) error {
	b.writes("((")
	// TODO: watch for passing an array type to writeCTypeName? In C, an array
	// type can decay into a pointer.
	if err := g.writeCTypeName(b, rhs, "", ""); err != nil {
		return err
	}
	b.writes(")(")
	if err := g.writeExpr(b, lhs, rp, parenthesesMandatory, depth); err != nil {
		return err
	}
	b.writes("))")
	return nil
}

func (g *gen) writeExprAssociativeOp(b *buffer, n *a.Expr, rp replacementPolicy, pp parenthesesPolicy, depth uint32) error {
	op := n.Operator()
	opName := cOpNames[0xFF&op]
	if opName == "" {
		return fmt.Errorf("unrecognized operator %q", op.AmbiguousForm().Str(g.tm))
	}

	if pp == parenthesesMandatory {
		b.writeb('(')
	}
	for i, o := range n.Args() {
		if i != 0 {
			b.writes(opName)
		}
		if err := g.writeExpr(b, o.Expr(), rp, parenthesesMandatory, depth); err != nil {
			return err
		}
	}
	if pp == parenthesesMandatory {
		b.writeb(')')
	}
	return nil
}

func (g *gen) writeCTypeName(b *buffer, n *a.TypeExpr, varNamePrefix string, varName string) error {
	// It may help to refer to http://unixwiz.net/techtips/reading-cdecl.html

	// TODO: fix this, allow slices of all types, not just of base.u8's. Also
	// allow arrays of slices, slices of pointers, etc.
	if n.IsSliceType() {
		o := n.Inner()
		if o.Decorator() == 0 && o.QID() == (t.QID{t.IDBase, t.IDU8}) && !o.IsRefined() {
			b.writes("wuffs_base__slice_u8")
			b.writeb(' ')
			b.writes(varNamePrefix)
			b.writes(varName)
			return nil
		}
		return fmt.Errorf("cannot convert Wuffs type %q to C", n.Str(g.tm))
	}

	// maxNumPointers is an arbitrary implementation restriction.
	const maxNumPointers = 16

	x := n
	for ; x != nil && x.IsArrayType(); x = x.Inner() {
	}

	numPointers, innermost := 0, x
	for ; innermost != nil && innermost.Inner() != nil; innermost = innermost.Inner() {
		// TODO: "nptr T", not just "ptr T".
		if p := innermost.Decorator(); p == t.IDPtr {
			if numPointers == maxNumPointers {
				return fmt.Errorf("cannot convert Wuffs type %q to C: too many ptr's", n.Str(g.tm))
			}
			numPointers++
			continue
		}
		// TODO: fix this.
		return fmt.Errorf("cannot convert Wuffs type %q to C", n.Str(g.tm))
	}

	fallback := true
	if qid := innermost.QID(); qid[0] == t.IDBase {
		if key := qid[1]; key < t.ID(len(cTypeNames)) {
			if s := cTypeNames[key]; s != "" {
				b.writes(s)
				fallback = false
			}
		}
	}
	if fallback {
		prefix := g.pkgPrefix
		qid := innermost.QID()
		if qid == (t.QID{t.IDBase, t.IDStatus}) {
			// No-op: special case "base.status" as being inside this package.
			//
			// TODO: change "base.status" in Wuffs code to just "status"? Or
			// change the C code's "wuffs_foo__status" to "wuffs_base__status"?
		} else if qid[0] != 0 {
			otherPkg := g.tm.ByID(qid[0])
			// TODO: map the "deflate" in "deflate.decoder" to the "deflate" in
			// `use "std/deflate"`, and use the latter "deflate".
			//
			// This is pretty academic at the moment, since they're the same
			// "deflate", but in the future, we might be able to rename used
			// packages, e.g. `use "foo/bar" as "baz"`, so "baz.qux" would map
			// to generating "wuffs_bar__qux".
			//
			// TODO: sanitize or validate otherPkg, e.g. that it's ASCII only?
			//
			// See gen.writeInitializerImpl for a similar use of otherPkg.
			prefix = "wuffs_" + otherPkg + "__"
		}
		b.printf("%s%s", prefix, qid[1].Str(g.tm))
	}

	for i := 0; i < numPointers; i++ {
		b.writeb('*')
	}

	b.writeb(' ')
	b.writes(varNamePrefix)
	b.writes(varName)

	x = n
	for ; x != nil && x.IsArrayType(); x = x.Inner() {
		b.writeb('[')
		b.writes(x.ArrayLength().ConstValue().String())
		b.writeb(']')
	}

	return nil
}

var cTypeNames = [...]string{
	t.IDI8:       "int8_t",
	t.IDI16:      "int16_t",
	t.IDI32:      "int32_t",
	t.IDI64:      "int64_t",
	t.IDU8:       "uint8_t",
	t.IDU16:      "uint16_t",
	t.IDU32:      "uint32_t",
	t.IDU64:      "uint64_t",
	t.IDBool:     "bool",
	t.IDIOReader: "wuffs_base__io_reader",
	t.IDIOWriter: "wuffs_base__io_writer",
}

var cOpNames = [256]string{
	t.IDEq:              " = ",
	t.IDPlusEq:          " += ",
	t.IDMinusEq:         " -= ",
	t.IDStarEq:          " *= ",
	t.IDSlashEq:         " /= ",
	t.IDShiftLEq:        " <<= ",
	t.IDShiftREq:        " >>= ",
	t.IDAmpEq:           " &= ",
	t.IDPipeEq:          " |= ",
	t.IDHatEq:           " ^= ",
	t.IDPercentEq:       " %= ",
	t.IDTildeModPlusEq:  " += ",
	t.IDTildeModMinusEq: " -= ",
	t.IDTildeSatPlusEq:  " no_such_C_operator ",
	t.IDTildeSatMinusEq: " no_such_C_operator ",

	t.IDXUnaryPlus:  " + ",
	t.IDXUnaryMinus: " - ",
	t.IDXUnaryNot:   " ! ",
	t.IDXUnaryRef:   " & ",
	t.IDXUnaryDeref: " * ",

	t.IDXBinaryPlus:          " + ",
	t.IDXBinaryMinus:         " - ",
	t.IDXBinaryStar:          " * ",
	t.IDXBinarySlash:         " / ",
	t.IDXBinaryShiftL:        " << ",
	t.IDXBinaryShiftR:        " >> ",
	t.IDXBinaryAmp:           " & ",
	t.IDXBinaryPipe:          " | ",
	t.IDXBinaryHat:           " ^ ",
	t.IDXBinaryPercent:       " % ",
	t.IDXBinaryTildeModPlus:  " + ",
	t.IDXBinaryTildeModMinus: " - ",
	t.IDXBinaryTildeSatPlus:  " no_such_C_operator ",
	t.IDXBinaryTildeSatMinus: " no_such_C_operator ",
	t.IDXBinaryNotEq:         " != ",
	t.IDXBinaryLessThan:      " < ",
	t.IDXBinaryLessEq:        " <= ",
	t.IDXBinaryEqEq:          " == ",
	t.IDXBinaryGreaterEq:     " >= ",
	t.IDXBinaryGreaterThan:   " > ",
	t.IDXBinaryAnd:           " && ",
	t.IDXBinaryOr:            " || ",
	t.IDXBinaryAs:            " no_such_C_operator ",

	t.IDXAssociativePlus: " + ",
	t.IDXAssociativeStar: " * ",
	t.IDXAssociativeAmp:  " & ",
	t.IDXAssociativePipe: " | ",
	t.IDXAssociativeHat:  " ^ ",
	t.IDXAssociativeAnd:  " && ",
	t.IDXAssociativeOr:   " || ",
}
