/*
Copyright 2023 The Go Authors. All rights reserved.
Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file.
*/

"use strict";

function fetchMessage() {
    var xmlHttp = new XMLHttpRequest();
    xmlHttp.open("GET", "/hello", false);
    xmlHttp.send(null);
    document.getElementById("message").innerHTML = xmlHttp.responseText;
}