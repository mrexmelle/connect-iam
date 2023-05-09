package organizationMember

import (
	"github.com/mrexmelle/connect-idp/internal/config"
)

type Repository struct {
	Config *config.Config
}

func NewRepository(cfg *config.Config) *Repository {
	return &Repository{
		Config: cfg,
	}
}

func (r *Repository) RetrieveByOrganizationIdWithKnownLeadEhid(
	organizationId string,
	leadEhid string,
) ([]Aggregate, error) {
	result, err := r.Config.Db.
		Table("profiles").
		Select("profiles.ehid, profiles.employee_id, profiles.name, profiles.email_address, tenures.title_name").
		Joins("LEFT JOIN tenures ON profiles.ehid = tenures.ehid").
		Where(
			"tenures.organization_id = ? AND tenures.start_date < NOW() AND (tenures.end_date IS NULL OR tenures.end_date > NOW())",
			organizationId,
		).
		Rows()
	if err != nil {
		return []Aggregate{}, err
	}
	defer result.Close()

	aggregates := make([]Aggregate, 0)
	for result.Next() {
		var agg = Aggregate{
			IsLead: false,
		}
		result.Scan(&agg.Ehid, &agg.EmployeeId, &agg.Name, &agg.EmailAddress, &agg.TitleName)
		aggregates = append(aggregates, agg)
	}
	return aggregates, nil
}
