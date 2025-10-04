package handler

import (
	"bytes"
	"html/template"

	"github.com/ia-edev-sindireceita/todo/internal/domain/application"
)

// TaskTemplateData holds data for rendering task HTML fragments
type TaskTemplateData struct {
	ID             string
	Title          string
	Description    string
	Status         string
	StatusClass    string
	StatusText     string
	CreatedAt      string
	ShowComplete   bool
	ShowShare      bool
	OwnershipClass string
	OwnershipText  string
}

var (
	// taskCardTemplate is the template for rendering a task card
	taskCardTemplate = template.Must(template.New("taskCard").Parse(`<div class="bg-white shadow rounded-lg p-6" id="task-{{.ID}}">
		<div class="flex justify-between items-start">
			<div class="flex-1">
				<h3 class="text-lg font-semibold text-gray-900">{{.Title}}</h3>
				<p class="text-gray-600 mt-1">{{.Description}}</p>
				<div class="mt-2 flex items-center space-x-2">
					<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium {{.StatusClass}}">
						{{.StatusText}}
					</span>
					<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium {{.OwnershipClass}}">
						{{.OwnershipText}}
					</span>
					<span class="text-sm text-gray-500">{{.CreatedAt}}</span>
				</div>
			</div>
			<div class="flex space-x-2 ml-4">
				{{if .ShowComplete}}
				<button hx-post="/web/tasks/{{.ID}}/complete" hx-target="#task-{{.ID}}" hx-swap="outerHTML"
						class="text-green-600 hover:text-green-800 font-medium">
					<svg class="w-4 h-4 inline mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
					</svg>
					Concluir
				</button>
				{{end}}
				{{if .ShowShare}}
				<button hx-post="/web/tasks/{{.ID}}/share"
						hx-target="#task-{{.ID}}"
						hx-swap="outerHTML"
						hx-prompt="Digite o email do usuário com quem deseja compartilhar:"
						hx-vals='js:{share_with_user_id: prompt("Digite o email do usuário:")}'
						class="text-blue-600 hover:text-blue-800 font-medium">
					<svg class="w-4 h-4 inline mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z"/>
					</svg>
					Compartilhar
				</button>
				{{end}}
				<button hx-delete="/web/tasks/{{.ID}}" hx-target="#task-{{.ID}}" hx-swap="outerHTML"
						hx-confirm="Tem certeza que deseja excluir esta tarefa?"
						class="text-red-600 hover:text-red-800">
					<svg class="w-4 h-4 inline mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
					</svg>
					Excluir
				</button>
			</div>
		</div>
	</div>`))

	// completedTaskTemplate is the template for rendering a completed task
	completedTaskTemplate = template.Must(template.New("completedTask").Parse(`<div class="bg-white shadow rounded-lg p-6" id="task-{{.ID}}">
		<div class="flex justify-between items-start">
			<div class="flex-1">
				<div class="flex items-center space-x-2">
					<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
						Concluída
					</span>
					<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium {{.OwnershipClass}}">
						{{.OwnershipText}}
					</span>
					<span class="text-sm text-gray-500">Tarefa concluída com sucesso!</span>
				</div>
			</div>
			<div class="flex space-x-2 ml-4">
				<button hx-delete="/web/tasks/{{.ID}}" hx-target="#task-{{.ID}}" hx-swap="outerHTML"
						hx-confirm="Tem certeza que deseja excluir esta tarefa?"
						class="text-red-600 hover:text-red-800">
					<svg class="w-4 h-4 inline mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
					</svg>
					Excluir
				</button>
			</div>
		</div>
	</div>`))
)

// renderTaskCard renders a task card HTML fragment with proper escaping
func renderTaskCard(task *application.Task, currentUserID string) (string, error) {
	isOwner := task.OwnerID == currentUserID

	data := TaskTemplateData{
		ID:           task.ID,
		Title:        task.Title,
		Description:  task.Description,
		Status:       string(task.Status),
		CreatedAt:    task.CreatedAt.Format("02/01/2006 15:04"),
		ShowComplete: task.Status == application.StatusPending,
		ShowShare:    isOwner && task.Status != application.StatusCompleted,
	}

	// Set status badge styling based on status
	switch task.Status {
	case application.StatusPending:
		data.StatusClass = "bg-yellow-100 text-yellow-800"
		data.StatusText = "Pendente"
	case application.StatusCompleted:
		data.StatusClass = "bg-green-100 text-green-800"
		data.StatusText = "Concluída"
	default:
		data.StatusClass = "bg-gray-100 text-gray-800"
		data.StatusText = string(task.Status)
	}

	// Set ownership badge styling based on owner
	if task.OwnerID == currentUserID {
		data.OwnershipClass = "bg-blue-100 text-blue-800"
		data.OwnershipText = "Própria"
	} else {
		data.OwnershipClass = "bg-purple-100 text-purple-800"
		data.OwnershipText = "Compartilhada"
	}

	var buf bytes.Buffer
	if err := taskCardTemplate.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// renderCompletedTask renders a completed task HTML fragment
func renderCompletedTask(task *application.Task, currentUserID string) (string, error) {
	data := TaskTemplateData{
		ID: task.ID,
	}

	// Set ownership badge styling based on owner
	if task.OwnerID == currentUserID {
		data.OwnershipClass = "bg-blue-100 text-blue-800"
		data.OwnershipText = "Própria"
	} else {
		data.OwnershipClass = "bg-purple-100 text-purple-800"
		data.OwnershipText = "Compartilhada"
	}

	var buf bytes.Buffer
	if err := completedTaskTemplate.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
