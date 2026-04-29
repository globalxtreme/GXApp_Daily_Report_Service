package database

type DatabaseSeeder interface {
	Seed()
}

type data struct {
	DatabaseSeeder
}

func Seeder() {
	seeders := []data{
		//{&seeder.TestingSeeder{}},
	}

	for _, seed := range seeders {
		seed.Seed()
	}
}
