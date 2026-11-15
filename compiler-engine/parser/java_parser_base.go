package parser

import "github.com/antlr4-go/antlr/v4"

type JavaParserBase struct {
	*antlr.BaseParser
}

func (p *JavaParserBase) IsVer15() bool {
	return true
}

func (p *JavaParserBase) IsVer16() bool {
	return true
}

func (p *JavaParserBase) IsNotIdentifierAssign() bool {
	// Dummy implementation for now
	return true
}

func (p *JavaParserBase) DoLastRecordComponent() bool {
	// Dummy implementation for now
	return true
}
