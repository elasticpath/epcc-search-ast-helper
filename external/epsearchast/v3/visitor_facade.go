package epsearchast_v3

type SearchFilterVisitor interface {
	PreVisit() error
	PostVisit() error
	PreVisitAnd() error
	PostVisitAnd() error
	VisitIn(args ...string) error
	VisitEq(first, second string) error
	VisitLe(first, second string) error
	VisitLt(first, second string) error
	VisitGe(first, second string) error
	VisitGt(first, second string) error
	VisitLike(first, second string) error
}

type SearchFilterVisitorAdaptor struct {
	Sfv SearchFilterVisitor
}

func (s *SearchFilterVisitorAdaptor) PreVisit() error {
	return s.Sfv.PreVisit()
}

func (s *SearchFilterVisitorAdaptor) PostVisit() error {
	return s.Sfv.PostVisit()
}

func (s *SearchFilterVisitorAdaptor) PreVisitAnd(astNode *AstNode) (bool, error) {
	return true, s.Sfv.PreVisitAnd()
}

func (s *SearchFilterVisitorAdaptor) PostVisitAnd(astNode *AstNode) (bool, error) {
	return true, s.Sfv.PostVisitAnd()
}

func (s *SearchFilterVisitorAdaptor) VisitIn(astNode *AstNode) (bool, error) {
	return false, s.Sfv.VisitIn(astNode.Args...)
}

func (s *SearchFilterVisitorAdaptor) VisitEq(astNode *AstNode) (bool, error) {
	return false, s.Sfv.VisitEq(astNode.FirstArg, astNode.SecondArg)
}

func (s *SearchFilterVisitorAdaptor) VisitLe(astNode *AstNode) (bool, error) {
	return false, s.Sfv.VisitLe(astNode.FirstArg, astNode.SecondArg)
}

func (s *SearchFilterVisitorAdaptor) VisitLt(astNode *AstNode) (bool, error) {
	return false, s.Sfv.VisitLt(astNode.FirstArg, astNode.SecondArg)
}

func (s *SearchFilterVisitorAdaptor) VisitGe(astNode *AstNode) (bool, error) {
	return false, s.Sfv.VisitGe(astNode.FirstArg, astNode.SecondArg)
}

func (s *SearchFilterVisitorAdaptor) VisitGt(astNode *AstNode) (bool, error) {
	return false, s.Sfv.VisitGt(astNode.FirstArg, astNode.SecondArg)
}

func (s *SearchFilterVisitorAdaptor) VisitLike(astNode *AstNode) (bool, error) {
	return false, s.Sfv.VisitLike(astNode.FirstArg, astNode.SecondArg)
}

var _ AstVisitor = (*SearchFilterVisitorAdaptor)(nil)
