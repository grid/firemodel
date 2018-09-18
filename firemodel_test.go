package firemodel_test

import (
	_ "github.com/mickeyreiss/firemodel/langs/go"
	_ "github.com/mickeyreiss/firemodel/langs/ios"
	_ "github.com/mickeyreiss/firemodel/langs/ts"

	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/mickeyreiss/firemodel"
	"github.com/sergi/go-diff/diffmatchpatch"
)

const fixturesRoot = "testfixtures/firemodel"

func TestFiremodelFromSchema(t *testing.T) {
	file, err := os.Open("firemodel.example.firemodel")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	schema, err := firemodel.ParseSchema(file)
	if err != nil {
		panic(err)
	}

	runTest(t, schema)
}

func (ctx *testCtx) firemodelConfig(testName string) *firemodel.Config {
	return &firemodel.Config{
		Languages: []firemodel.Language{
			{Language: "ios", Output: "./swift/"},
			{Language: "go", Output: "./go"},
			{Language: "ts", Output: "./ts/"},
		},
		SourceCoderProvider: ctx.newTestSourceCodeProvider(testName),
	}
}

type inMemoryFilesByName map[string]*inMemoryFile

func (filesByName inMemoryFilesByName) keys() []string {
	ret := make([]string, 0, len(filesByName))
	for name := range filesByName {
		ret = append(ret, name)
	}
	return ret
}

type testCtx struct {
	prefix string
	files  inMemoryFilesByName
}

func runTest(t *testing.T, schema *firemodel.Schema) {
	ctx := testCtx{
		prefix: path.Join(fixturesRoot, t.Name()),
		files:  make(inMemoryFilesByName),
	}

	if isUpdate() {
		if err := os.RemoveAll(ctx.prefix); err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(ctx.prefix, 0700); err != nil {
			t.Fatal(err)
		}
	}

	if err := firemodel.Run(schema, ctx.firemodelConfig(t.Name())); err != nil {
		t.Fatal(err)
	}

	if !isUpdate() {
		fixtures, err := filepath.Glob(path.Join(ctx.prefix, "*", "*"))
		if err != nil {
			panic(err)
		}
		if len(fixtures) != len(ctx.files) {
			t.Errorf("Fixtures do not match up with generated files. Fixtures %v, Actual: %v", fixtures, ctx.files.keys())
		}
		for _, filename := range fixtures {
			actualBuf, ok := ctx.files[filename]
			if !ok {
				t.Errorf("Missing generated file for fixture %s", filename)
				continue
			}
			actual := actualBuf.Bytes()
			exp, err := ioutil.ReadFile(filename)
			if err != nil {
				t.Fatalf("fixture not found for generated file: %v", err)
			}
			if !bytes.Equal(exp, actual) {
				dmp := diffmatchpatch.New()
				diffs := dmp.DiffMain(string(exp), string(actual), true)

				t.Error(dmp.DiffPrettyText(diffs))
				t.Log("If this diff looks ok, re-run tests with FIREMODEL_UPDATE_FIXTURES=true")
			}
		}
	}
}

func (ctx *testCtx) newTestSourceCodeProvider(testName string) func(prefix string) firemodel.SourceCoder {
	return func(prefix string) firemodel.SourceCoder {
		coder := &testSourceCoder{
			prefix: path.Join(ctx.prefix, prefix),
			files:  ctx.files,
		}
		if err := os.MkdirAll(coder.prefix, 0700); err != nil {
			panic(err)
		}
		return coder
	}
}

type testSourceCoder struct {
	prefix string
	files  map[string]*inMemoryFile
}

func isUpdate() bool {
	_, update := os.LookupEnv("FIREMODEL_UPDATE_FIXTURES")
	return update
}

func (ctx *testSourceCoder) NewFile(filename string) (io.WriteCloser, error) {
	if isUpdate() {
		return os.OpenFile(path.Join(ctx.prefix, filename), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	} else {
		var file inMemoryFile
		ctx.files[path.Join(ctx.prefix, filename)] = &file
		return &file, nil
	}
}

func (ctx *testSourceCoder) Flush() error {
	return nil
}

type inMemoryFile struct {
	bytes.Buffer
}

func (f *inMemoryFile) Write(p []byte) (n int, err error) {
	return f.Buffer.Write(p)
}

func (_ *inMemoryFile) Close() error {
	return nil
}
