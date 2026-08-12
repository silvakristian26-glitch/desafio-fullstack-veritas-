import { useEffect, useState, useCallback } from 'react'
import { api } from './api'
import { COLUMNS } from './constants'
import Column from './components/Column'
import TaskModal from './components/TaskModal'

export default function App() {
  const [tasks, setTasks] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [modal, setModal] = useState(null)
  const [saving, setSaving] = useState(false)
  const [busyId, setBusyId] = useState(null)

  const loadTasks = useCallback(async () => {
    setLoading(true)
    setError('')
    try {
      const data = await api.listTasks()
      setTasks(data || [])
    } catch (err) {
      setError(err.message || 'Não foi possível carregar as tarefas.')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    loadTasks()
  }, [loadTasks])

  async function handleSave(payload) {
    setSaving(true)
    setError('')
    try {
      if (modal.mode === 'edit') {
        const updated = await api.updateTask(modal.task.id, payload)
        setTasks((prev) => prev.map((t) => (t.id === updated.id ? updated : t)))
      } else {
        const created = await api.createTask(payload)
        setTasks((prev) => [...prev, created])
      }
      setModal(null)
    } catch (err) {
      setError(err.message || 'Não foi possível salvar a tarefa.')
    } finally {
      setSaving(false)
    }
  }

  async function handleDelete(task) {
    if (!window.confirm(`Excluir a tarefa "${task.title}"?`)) return
    setBusyId(task.id)
    setError('')
    try {
      await api.deleteTask(task.id)
      setTasks((prev) => prev.filter((t) => t.id !== task.id))
    } catch (err) {
      setError(err.message || 'Não foi possível excluir a tarefa.')
    } finally {
      setBusyId(null)
    }
  }

  async function handleMove(taskId, newStatus) {
    const task = tasks.find((t) => t.id === taskId)
    if (!task || task.status === newStatus) return

    const previous = tasks
    setTasks((prev) =>
      prev.map((t) => (t.id === taskId ? { ...t, status: newStatus } : t)),
    )
    setBusyId(taskId)
    setError('')
    try {
      const updated = await api.updateTask(taskId, {
        title: task.title,
        description: task.description,
        status: newStatus,
      })
      setTasks((prev) => prev.map((t) => (t.id === updated.id ? updated : t)))
    } catch (err) {
      setTasks(previous)
      setError(err.message || 'Não foi possível mover a tarefa.')
    } finally {
      setBusyId(null)
    }
  }

  return (
    <div className="app">
      <header className="app__header">
        <div>
          <p className="app__eyebrow">Painel de Tarefas</p>
          <h1>Kanban</h1>
        </div>
        <button type="button" className="btn btn--primary" onClick={() => setModal({ mode: 'create' })}>
          + Nova tarefa
        </button>
      </header>

      {error && (
        <div className="banner banner--error" role="alert">
          {error}
          <button type="button" onClick={() => setError('')} aria-label="Fechar aviso">
            ✕
          </button>
        </div>
      )}

      {loading ? (
        <p className="app__loading">Carregando tarefas…</p>
      ) : (
        <main className="board">
          {COLUMNS.map((column) => (
            <Column
              key={column.status}
              column={column}
              tasks={tasks.filter((t) => t.status === column.status)}
              busyId={busyId}
              onAdd={(status) => setModal({ mode: 'create', defaultStatus: status })}
              onEdit={(task) => setModal({ mode: 'edit', task })}
              onDelete={handleDelete}
              onMove={handleMove}
            />
          ))}
        </main>
      )}

      {modal && (
        <TaskModal
          initial={
            modal.mode === 'edit'
              ? modal.task
              : modal.defaultStatus
              ? { status: modal.defaultStatus }
              : null
          }
          saving={saving}
          onSave={handleSave}
          onClose={() => setModal(null)}
        />
      )}
    </div>
  )
}
