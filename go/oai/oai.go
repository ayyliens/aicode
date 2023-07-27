package oai

const TempDirName = `aicode`

type OaiFunction interface{ OaiCall(string) string }
