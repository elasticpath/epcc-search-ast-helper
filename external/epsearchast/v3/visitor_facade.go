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

func NewSearchFilterVisitorAdapter(visitor SearchFilterVisitor) AstVisitor {
	return &SearchFilterVisitorAdaptor{Sfv: visitor}
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

func (s *SearchFilterVisitorAdaptor) PreVisitAnd(_ *AstNode) (bool, error) {
	return true, s.Sfv.PreVisitAnd()
}

func (s *SearchFilterVisitorAdaptor) PostVisitAnd(_ *AstNode) error {
	return s.Sfv.PostVisitAnd()
}

func (s *SearchFilterVisitorAdaptor) VisitIn(astNode *AstNode) (bool, error) {
	return false, s.Sfv.VisitIn(astNode.Args...)
}

func (s *SearchFilterVisitorAdaptor) VisitEq(astNode *AstNode) (bool, error) {
	return false, s.Sfv.VisitEq(astNode.Args[0], astNode.Args[1])
}

func (s *SearchFilterVisitorAdaptor) VisitLe(astNode *AstNode) (bool, error) {
	return false, s.Sfv.VisitLe(astNode.Args[0], astNode.Args[1])
}

func (s *SearchFilterVisitorAdaptor) VisitLt(astNode *AstNode) (bool, error) {
	return false, s.Sfv.VisitLt(astNode.Args[0], astNode.Args[1])
}

func (s *SearchFilterVisitorAdaptor) VisitGe(astNode *AstNode) (bool, error) {
	return false, s.Sfv.VisitGe(astNode.Args[0], astNode.Args[1])
}

func (s *SearchFilterVisitorAdaptor) VisitGt(astNode *AstNode) (bool, error) {
	return false, s.Sfv.VisitGt(astNode.Args[0], astNode.Args[1])
}

func (s *SearchFilterVisitorAdaptor) VisitLike(astNode *AstNode) (bool, error) {
	return false, s.Sfv.VisitLike(astNode.Args[0], astNode.Args[1])
}

var _ AstVisitor = (*SearchFilterVisitorAdaptor)(nil)
