<h1>linx</h1>  
<p> Reveals invisible links within JavaScript files. Inspired by <a href="https://github.com/GerbenJavado/LinkFinder">LinkFinder</a> </p>  
<p>  
  <a href="https://opensource.org/licenses/MIT">  
    <img src="https://img.shields.io/badge/license-MIT-_red.svg">  
  </a>  
  <a href="https://goreportcard.com/badge/github.com/riza/linx">  
    <img src="https://goreportcard.com/badge/github.com/riza/linx">  
  </a>  
  <a href="https://github.com/riza/linx/releases">  
    <img src="https://img.shields.io/github/release/riza/linx">  
  </a>  
  <a href="https://twitter.com/rizasabuncu">  
    <img src="https://img.shields.io/twitter/follow/rizasabuncu.svg?logo=twitter">  
  </a>  
</p>

# Installation

linx requires **go1.17** to install successfully. Run the following command to get the repo -

```sh
go install -v github.com/riza/linx/cmd/linx@latest
```

# Usage

```sh
linx --target=https://rizasabuncu.com/assets/admin_acces.js
```

# TODOs

* [x] HTML output support 
* [ ] JSON output support
* [ ] Custom cookie support
* [ ] Rule improvement & blacklist support
* [ ] Support parallel scan multiple files
* [ ] ...
