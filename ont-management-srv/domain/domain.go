package domain

type DomainRepository interface {
	FindDomainCategory(domain string) (int, error)
}
