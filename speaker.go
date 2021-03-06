// Copyright ©2017 The ev3go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package evb

import "github.com/ev3go/ev3dev"

// Speaker is a handle to the evb speaker. It must be initialized before use.
var Speaker = ev3dev.NewSpeaker("/dev/input/by-path/platform-evb-sound-event")
