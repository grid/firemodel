package firemodel_test

import (
	_ "github.com/mickeyreiss/firemodel/langs/ios"
	_ "github.com/mickeyreiss/firemodel/langs/go"
	_ "github.com/mickeyreiss/firemodel/langs/ts"

	"testing"
	"github.com/spf13/viper"
	"github.com/pkg/errors"
	"os"
	"io"
	"io/ioutil"
	"go.uber.org/zap/buffer"
	"fmt"
	"strings"
	"path"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/mickeyreiss/firemodel"
)

const fixturesRoot = "testfixtures"

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

	ctx := &testCtx{t: t}

	if err := firemodel.Run(schema, ctx.firemodelConfig()); err != nil {
		panic(err)
	}

	ctx.Test()
}

func TestFiremodelFromYamlSpec(t *testing.T) {
	// Set up config.
	v := viper.New()
	v.SetConfigName("firemodel.example")
	v.AddConfigPath(".")
	if err := v.ReadInConfig(); err != nil {
		panic(errors.Wrap(err, "firemodel: reading config"))
	}
	var schema firemodel.Schema
	if err := v.UnmarshalExact(&schema); err != nil {
		panic(errors.Wrapf(err, "firemodel: parsing %s", v.ConfigFileUsed()))
	}

	// Set up output.
	ctx := &testCtx{t: t}

	if err := firemodel.Run(&schema, ctx.firemodelConfig()); err != nil {
		panic(err)
	}

	ctx.Test()
}

func (ctx *testCtx) firemodelConfig() *firemodel.Config {
	ctx.t.Helper()

	return &firemodel.Config{
		Languages: []firemodel.Language{
			{Language: "ios", Output: "./swift/"},
			{Language: "go", Output: "./go"},
			{Language: "ts", Output: "./ts/"},
		},
		SourceCoderProvider: func(prefix string) firemodel.SourceCoder {
			fmt.Fprintf(&ctx.buf, "===================== %s ===================\n", prefix)
			return ctx
		},
	}
}

type testCtx struct {
	buf buffer.Buffer
	t   *testing.T
}

func (ctx *testCtx) Test() {
	if _, update := os.LookupEnv("FIREMODEL_UPDATE_FIXTURES"); update {
		os.MkdirAll(fixturesRoot, 0700)
		err := ioutil.WriteFile(ctx.fixtureFile(), ctx.buf.Bytes(), 0600)
		if err != nil {
			panic(err)
		}
	} else {
		exp, err := ioutil.ReadFile(ctx.fixtureFile())
		if err != nil {
			ctx.t.Error(err)
		}

		res := ctx.buf.String()
		if strings.TrimSpace(string(exp)) != strings.TrimSpace(res) {

			dmp := diffmatchpatch.New()
			diffs := dmp.DiffMain(string(exp), res, true)

			ctx.t.Error(dmp.DiffPrettyText(diffs))
			ctx.t.Log("If this diff looks ok, re-run tests with FIREMODEL_UPDATE_FIXTURES=true")
		}
	}
}

func (ctx *testCtx) fixtureFile() string {
	return path.Join(fixturesRoot, ctx.t.Name()+".txt")
}

func (ctx *testCtx) NewFile(filename string) (io.WriteCloser, error) {
	fmt.Fprintf(&ctx.buf, "===================== Open %s ===================\n", filename)
	return ctx, nil
}

func (ctx *testCtx) Write(p []byte) (n int, err error) {
	return ctx.buf.Write(p)
}

func (ctx *testCtx) Close() error {
	_, err := fmt.Fprintf(&ctx.buf, "===================== Close ===================\n")
	return err
}

func (ctx *testCtx) Flush() error {
	_, err := fmt.Fprintf(&ctx.buf, "===================== Flush ===================\n")
	return err
}
