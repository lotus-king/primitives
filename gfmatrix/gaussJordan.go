package gfmatrix

// gaussJordan reduces the matrix according to the Gauss-Jordan Method.  Returns the augment matrix, the transformed
// matrix, and the set set of free variables.
func (e Matrix) gaussJordan() (aug, f Matrix, frees []int) {
	out, in := e.Size()

	aug = GenerateIdentity(in)
	for x := in; x < out; x++ {
		aug = append(aug, NewRow(in))
	}

	f = e.Dup() // Duplicate e away so we don't mutate it.

	row, col := 0, 0

	for row < out && col < in {
		// Find a non-zero element to move into the pivot position.
		i := f.FindPivot(row, col)
		if i == -1 { // Failed to find a pivot.
			frees = append(frees, col)
			col++

			continue
		}

		// Move it into position.
		f[row], f[i] = f[i], f[row]
		aug[row], aug[i] = aug[i], aug[row]

		// Normalize the entry in this rows (pivot)th position.
		correction := f[row][col].Invert()
		f[row] = f[row].ScalarMul(correction)
		aug[row] = aug[row].ScalarMul(correction)

		// Cancel out the (row)th position for every row above and below it.
		for j, _ := range f {
			if j != row && !f[j][col].IsZero() {
				correction := f[j][col].Dup()
				f[j] = f[j].Add(f[row].ScalarMul(correction))
				aug[j] = aug[j].Add(aug[row].ScalarMul(correction))
			}
		}

		row++
		col++
	}

	// Add the rest of the free variables for completion.
	for free := out; free < in; free++ {
		frees = append(frees, free)
	}

	return
}

// NullSpace returns a basis for the matrix's nullspace.
func (e Matrix) NullSpace() (basis []Row) {
	out, in := e.Size()
	if out == 0 {
		return []Row{}
	}

	_, f, frees := e.gaussJordan()

	for _, free := range frees {
		input := NewRow(in)
		input[free] = 0x01

		for _, row := range f[:in] {
			if !row[free].IsZero() {
				input[row.Height()] = row[free]
			}
		}

		basis = append(basis, input)
	}

	return
}
