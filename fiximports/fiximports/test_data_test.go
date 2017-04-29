// Copyright 2017 Corey Scott http://www.sage42.org/
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fiximports

var testFile1 = `package main

import (
	"github.com/corsc/go-tools/commons"

	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {}
`

var testFile1Fixed = `package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/corsc/go-tools/commons"
)

func main() {}
`

var testFile2 = `package main

import (
	"net/http/httptest"
	"net/http"
	"net/http/httputil"
)

func main() {}
`

var testFile2Fixed = `package main

import (
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
)

func main() {}
`

var testFile3 = `package main

import (
	"crypto/hmac"
	"crypto/subtle"
	"fmt"
	"math"
	"time"
	"crypto/sha256"
	"encoding/base64"
)

func main() {}
`

var testFile3Fixed = `package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"math"
	"time"
)

func main() {}
`

var testFileNoImports = `package main

func main() {}
`

var testFileDotImport = `package main

import (
	"fmt"
	. "net/http/httputil"
)

func main() {}
`

var testFileBlankImport = `package main

import (
	"fmt"
	_ "net/http/httputil"
)

func main() {}
`

var testFileCommentedImportAbove = `package main

import (
	"fmt"
	// comment above
	"net/http/httputil"
)

func main() {}
`

var testFileCommentedImportAtEnd = `package main

import (
	"fmt"
	"net/http/httputil" // comment on the end
)

func main() {}
`

var testFileIndividualImports = `package main

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"

import io "io"

func main() {}
`

var testFileIndividualImportsFixed = `package main

import (
	"fmt"
	"io"
	"math"

	_ "github.com/gogo/protobuf/gogoproto"
	"github.com/golang/protobuf/proto"
)

func main() {}
`

var testFileSingleImport = `package main

import "github.com/golang/protobuf/proto"

func main() {}
`

var testFileExtraLine = `package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"math"
	"os"
	"regexp"
	"sort"

	"strings"

	"github.com/corsc/go-tools/commons"
)

func main() {}
`

var testFileExtraLineFixed = `package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"math"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/corsc/go-tools/commons"
)

func main() {}
`
