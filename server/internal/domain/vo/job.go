package vo

import "fmt"

type Job struct {
	JobID               string
	JobType             JobType
	AvailableWeaponList []WeaponType
}

func NewJob(id string, jobTypeStr string, availableWeaponList []WeaponType) (Job, error) {
	jobType := ToJobType(jobTypeStr)
	if jobType == JobTypeUnspecified {
		return Job{}, fmt.Errorf("job type is unspecified: %s", jobTypeStr)
	}
	return Job{
		JobID:               id,
		JobType:             jobType,
		AvailableWeaponList: availableWeaponList,
	}, nil
}

func (j Job) ID() string {
	return j.JobID
}

type JobType string

func (j JobType) String() string {
	return string(j)
}

func ToJobType(s string) JobType {
	switch s {
	case JobTypeWarrior.String():
		return JobTypeWarrior
	case JobTypeKnight.String():
		return JobTypeKnight
	case JobTypeDarkKnight.String():
		return JobTypeDarkKnight
	case JobTypeThief.String():
		return JobTypeThief
	case JobTypeAssassin.String():
		return JobTypeAssassin
	case JobTypeHermit.String():
		return JobTypeHermit
	default:
		return JobTypeUnspecified
	}
}

const (
	JobTypeUnspecified JobType = "unspecified"
	JobTypeWarrior     JobType = "warrior"
	JobTypeKnight      JobType = "knight"
	JobTypeDarkKnight  JobType = "dark knight"
	JobTypeThief       JobType = "thief"
	JobTypeAssassin    JobType = "assassin"
	JobTypeHermit      JobType = "hermit"
)

//var jobAvailableWeaponTypes = map[JobType][]WeaponType{
//	JobTypeWarrior:    {WeaponTypeKnife, WeaponTypeShield},
//	JobTypeKnight:     {WeaponTypeKnife, WeaponTypeShield},
//	JobTypeDarkKnight: {WeaponTypeKnife, WeaponTypeShield},
//	JobTypeThief:      {WeaponTypeDagger, WeaponTypeShield},
//	JobTypeAssassin:   {WeaponTypeKnife, WeaponTypeDagger, WeaponTypeShield},
//	JobTypeHermit:     {WeaponTypeKnife, WeaponTypeDagger, WeaponTypeShield},
//}
