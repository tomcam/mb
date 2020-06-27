package mdext

// func TestInferTitle(t *testing.T) {
// 	type args struct {
// 		root  ast.Node
// 		mdSrc []byte
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want string
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := InferTitle(tt.args.root, tt.args.mdSrc); got != tt.want {
// 				t.Errorf("InferTitle() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
//
// func TestInferTitle(t *testing.T) {
// 	tests := []struct {
// 		mdSrc string
// 		want  string
// 	}{
// 		{texts.Dedent(`
// 		     # h1.1
// 		     body
// 		     ## h2.1
// 		`), "h1.1"},
// 		{texts.Dedent(`
// 		     # h1.1 *foo* bar
// 		     body
// 		     ## h2.1
// 		`), "h1.1 foo bar"},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.mdSrc, func(t *testing.T) {
// 			app := newApp(tt.mdSrc)
// 			app.Site.MarkdownOptions.headingIDs = true
// 			tocs := app.generateTOC(tt.level)
// 			if diff := cmp.Diff(tt.want, tocs); diff != "" {
// 				t.Errorf("extractTOCs() mismatch (-want +got):\n%s", diff)
// 			}
// 		})
// 	}
// }
