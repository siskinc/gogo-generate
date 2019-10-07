//go:generate ./../generate-mgorm --type MogormTest
package test

//@def MogormTest struct comment1
//@def MogormTest struct comment2
//@def MogormTest struct comment3
type MogormTest struct {
	//das
	AA int //@def aa
	BB int //@def bb
}

/*@def MogormTest struct comment1
@def MogormTest struct comment2
@def MogormTest struct comment3*/
type MogormTest2 struct {
	//das
	AA int //@def aa
}
