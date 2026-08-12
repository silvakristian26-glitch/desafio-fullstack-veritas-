import { COLUMNS, COLUMN_BY_STATUS } from '../constants'

function formatDate(iso) {
  try {
    return new Date(iso).toLocaleDateString('pt-BR', {
      day: '2-digit',
      month: '2-digit',
    })
  } catch {
    return ''
  }
}

export default function TaskCard({ task, onEdit, onDelete, onMove, busy }) {
  const accent = COLUMN_BY_STATUS[task.status]?.accent || '#999'

  return (
    <article
      className="task-card"
      style={{ borderLeftColor: accent }}
      draggable
      onDragStart={(e) => {
        e.dataTransfer.setData('text/task-id', task.id)
      }}
    >
      <header className="task-card__header">
        <h3>{task.title}</h3>
        <span className="task-card__id">#{task.id.slice(0, 6)}</span>
      </header>

      {task.description && <p className="task-card__desc">{task.description}</p>}

      <footer className="task-card__footer">
        <span className="task-card__date">{formatDate(task.updated_at)}</span>

        <div className="task-card__actions">
          <select
            aria-label="Mover para outra coluna"
            value={task.status}
            disabled={busy}
            onChange={(e) => onMove(task, e.target.value)}
          >
            {COLUMNS.map((c) => (
              <option key={c.status} value={c.status}>
                {c.label}
              </option>
            ))}
          </select>
          <button
            type="button"
            className="icon-btn"
            title="Editar"
            disabled={busy}
            onClick={() => onEdit(task)}
          >
            ✎
          </button>
          <button
            type="button"
            className="icon-btn icon-btn--danger"
            title="Excluir"
            disabled={busy}
            onClick={() => onDelete(task)}
          >
            ✕
          </button>
        </div>
      </footer>
    </article>
  )
}
