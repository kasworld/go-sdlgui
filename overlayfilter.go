// Copyright 2015 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sdlgui

import (
	"github.com/kasworld/htmlcolors"
)

type OverlayFilter []htmlcolors.RGBA

func MakeOverlayFilter(n int, b htmlcolors.RGBA) OverlayFilter {
	rtn := make(OverlayFilter, n)
	for i := 0; i < n; i++ {
		for j, v := range b {
			if int(v) < i*256/n {
				rtn[i][j] = 0
			} else {
				rtn[i][j] = v - uint8(i*256/n)
			}
		}
	}
	return rtn
}
