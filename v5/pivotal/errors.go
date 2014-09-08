/*
   Copyright (C) 2014  Salsita s.r.o.

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program. If not, see {http://www.gnu.org/licenses/}.
*/

package pivotal

import (
	"fmt"
	"net/http"
)

type ErrFieldNotSet struct {
	fieldName string
}

func (err *ErrFieldNotSet) Error() string {
	return fmt.Sprintf("Required field '%s' is not set")
}

type ErrHTTP struct {
	*http.Response
}

func (err *ErrHTTP) Error() string {
	return fmt.Sprintf("%v %v -> %v", err.Request.Method, err.Request.URL, err.Status)
}
