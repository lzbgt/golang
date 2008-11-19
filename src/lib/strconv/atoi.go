// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strconv
import "os"

func IntSize() uint {
	siz := uint(8);
	for 1<<siz != 0 {
		siz *= 2
	}
	return siz
}
var intsize = IntSize();

// Convert decimal string to unsigned integer.
export func atoui64(s string) (i uint64, err *os.Error) {
	// empty string bad
	if len(s) == 0 {
		return 0, os.EINVAL
	}

	// pick off zero
	if s == "0" {
		return 0, nil
	}

	// otherwise, leading zero bad:
	// don't want to take something intended as octal.
	if s[0] == '0' {
		return 0, os.EINVAL
	}

	// parse number
	n := uint64(0);
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return 0, os.EINVAL
		}
		if n > (1<<64)/10 {
			return 1<<64-1, os.ERANGE
		}
		n = n*10;
		d := uint64(s[i] - '0');
		if n+d < n {
			return 1<<64-1, os.ERANGE
		}
		n += d;
	}
	return n, nil
}

// Convert decimal string to integer.
export func atoi64(s string) (i int64, err *os.Error) {
	// empty string bad
	if len(s) == 0 {
		return 0, os.EINVAL
	}

	// pick off leading sign
	neg := false;
	if s[0] == '+' {
		s = s[1:len(s)]
	} else if s[0] == '-' {
		neg = true;
		s = s[1:len(s)]
	}

	var un uint64;
	un, err = atoui64(s);
	if err != nil && err != os.ERANGE {
		return 0, err
	}
	if !neg && un >= 1<<63 {
		return 1<<63-1, os.ERANGE
	}
	if neg && un > 1<<63 {
		return -1<<63, os.ERANGE
	}
	n := int64(un);
	if neg {
		n = -n
	}
	return n, nil
}

export func atoui(s string) (i uint, err *os.Error) {
	i1, e1 := atoui64(s);
	if e1 != nil && e1 != os.ERANGE {
		return 0, e1
	}
	i = uint(i1);
	if uint64(i) != i1 {
		// TODO: return uint(^0), os.ERANGE.
		i1 = 1<<64-1;
		return uint(i1), os.ERANGE
	}
	return i, nil
}

export func atoi(s string) (i int, err *os.Error) {
	i1, e1 := atoi64(s);
	if e1 != nil && e1 != os.ERANGE {
		return 0, e1
	}
	i = int(i1);
	if int64(i) != i1 {
		if i1 < 0 {
			return -1<<(intsize-1), os.ERANGE
		}
		return 1<<(intsize-1) - 1, os.ERANGE
	}
	return i, nil
}

