export const COLUMNS = [
  { status: 'a_fazer', label: 'A Fazer', accent: '#e3a008' },
  { status: 'em_progresso', label: 'Em Progresso', accent: '#3e7cb1' },
  { status: 'concluida', label: 'Concluídas', accent: '#4c9a6a' },
]

export const COLUMN_BY_STATUS = Object.fromEntries(
  COLUMNS.map((c) => [c.status, c]),
)
