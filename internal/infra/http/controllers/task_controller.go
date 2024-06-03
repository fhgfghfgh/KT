package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/app"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/requests"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/resources"
	"github.com/go-chi/chi/v5"
)

type TaskController struct {
	taskService app.TaskService
}

func NewTaskController(ts app.TaskService) TaskController {
	return TaskController{
		taskService: ts,
	}
}

func (c TaskController) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", c.Save())
	r.Get("/user/{userId}", c.FindByUserId())
	r.Get("/{taskId}", c.FindById())
	r.Put("/{taskId}", c.Update())
	r.Delete("/{taskId}", c.Delete())

	return r
}

func (c TaskController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		task, err := requests.Bind(r, requests.TaskRequest{}, domain.Task{})
		if err != nil {
			log.Printf("TaskController -> Save: %s", err)
			BadRequest(w, err)
			return
		}

		user := r.Context().Value(UserKey).(domain.User)
		task.UserId = user.Id
		task.Status = domain.NewTaskStatus
		task, err = c.taskService.Save(task)
		if err != nil {
			log.Printf("TaskController -> Save: %s", err)
			InternalServerError(w, err)
			return
		}

		var tDto resources.TaskDto
		tDto = tDto.DomainToDto(task)
		Created(w, tDto)
	}
}

func (c TaskController) FindByUserId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		tasks, err := c.taskService.FindByUserId(user.Id)
		if err != nil {
			log.Printf("TaskController -> FindByUserId: %s", err)
			InternalServerError(w, err)
			return
		}

		var tsDto resources.TasksDto
		tsDto = tsDto.DomainToDtoCollection(tasks)
		Success(w, tsDto)
	}
}

func (c TaskController) FindById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskIdStr := chi.URLParam(r, "taskId")
		taskId, err := strconv.ParseUint(taskIdStr, 10, 64)
		if err != nil {
			log.Printf("TaskController -> FindById: %s", err)
			BadRequest(w, err)
			return
		}

		task, err := c.taskService.FindById(taskId)
		if err != nil {
			log.Printf("TaskController -> FindById: %s", err)
			InternalServerError(w, err)
			return
		}

		var tDto resources.TaskDto
		tDto = tDto.DomainToDto(task)
		Success(w, tDto)
	}
}

func (c TaskController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskIdStr := chi.URLParam(r, "taskId")
		taskId, err := strconv.ParseUint(taskIdStr, 10, 64)
		if err != nil {
			log.Printf("TaskController -> Update: %s", err)
			BadRequest(w, err)
			return
		}

		task, err := requests.Bind(r, requests.TaskRequest{}, domain.Task{})
		if err != nil {
			log.Printf("TaskController -> Update: %s", err)
			BadRequest(w, err)
			return
		}

		task.Id = taskId
		task, err = c.taskService.Update(task)
		if err != nil {
			log.Printf("TaskController -> Update: %s", err)
			InternalServerError(w, err)
			return
		}

		var tDto resources.TaskDto
		tDto = tDto.DomainToDto(task)
		Success(w, tDto)
	}
}

func (c TaskController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskIdStr := chi.URLParam(r, "taskId")
		taskId, err := strconv.ParseUint(taskIdStr, 10, 64)
		if err != nil {
			log.Printf("TaskController -> Delete: %s", err)
			BadRequest(w, err)
			return
		}

		err = c.taskService.Delete(taskId)
		if err != nil {
			log.Printf("TaskController -> Delete: %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, nil)
	}
}
