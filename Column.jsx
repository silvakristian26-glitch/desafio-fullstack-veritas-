import { useState } from 'react'
import TaskCard from './TaskCard'

export default function Column({ column, tasks, onEdit, onDelete, onMove, onAdd, busyId }) {
  const [dragOver, setDragOver] = useState(false)

  return (
    <section
      className={`column${dragOver ? ' column--drag-over' : ''}`}
      style={{ '--accent': column.accent }}
      onDragOver={(e) => {
        e.preventDefault()
        setDragOver(true)
      }}
      onDragLeave={() => setDragOver(false)}
      onDrop={(e) => {
        e.preventDefault()
        setDragOver(false)
        const taskId = e.dataTransfer.getData('text/task-id')
        const task = tasks.find((t) => t.id === taskId)
        if (task) return
        onMove(taskId, column.status)
      }}
    >
      <header className="column__header">
        <span className="column__dot" />
        <h2>{column.label}</h2>
        <span className="column__count">{tasks.length}</span>
      </header>

      <button type="button" className="column__add" onClick={() => onAdd(column.status)}>
        + Nova tarefa
      </button>

      <div className="column__list">
        {tasks.length === 0 && (
          <p className="column__empty">Nenhuma tarefa aqui ainda.</p>
        )}
        {tasks.map((task) => (
          <TaskCard
            key={task.id}
            task={task}
            onEdit={onEdit}
            onDelete={onDelete}
            onMove={(t, status) => onMove(t.id, status)}
            busy={busyId === task.id}
          />
        ))}
      </div>
    </section>
  )
}
