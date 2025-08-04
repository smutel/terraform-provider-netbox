# NOTE: This file is for HashiCorp specific licensing automation and can be deleted after creating a new repo with this template.
schema_version = 1

project {
  license        = "MIT"
  copyright_holder = ""
  copyright_year = 2019

  header_ignore = [
    # internal catalog metadata (prose)
    "META.d/**/*.yaml",

    # examples used within documentation (prose)
    "examples/**",

    # GitHub issue template configuration
    ".github/ISSUE_TEMPLATE/*.yml",

    # golangci-lint tooling configuration
    ".golangci.yml",

    # GoReleaser tooling configuration
    ".goreleaser.yml",
  ]
}
