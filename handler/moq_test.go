// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package handler

import (
	"context"
	"github.com/YutaKakiki/go-todo-api/entity"
	"sync"
)

// Ensure, that ListTasksServiceMock does implement ListTasksService.
// If this is not the case, regenerate this file with moq.
var _ ListTasksService = &ListTasksServiceMock{}

// ListTasksServiceMock is a mock implementation of ListTasksService.
//
//	func TestSomethingThatUsesListTasksService(t *testing.T) {
//
//		// make and configure a mocked ListTasksService
//		mockedListTasksService := &ListTasksServiceMock{
//			ListTaskFunc: func(ctx context.Context) (entity.Tasks, error) {
//				panic("mock out the ListTask method")
//			},
//		}
//
//		// use mockedListTasksService in code that requires ListTasksService
//		// and then make assertions.
//
//	}
type ListTasksServiceMock struct {
	// ListTaskFunc mocks the ListTask method.
	ListTaskFunc func(ctx context.Context) (entity.Tasks, error)

	// calls tracks calls to the methods.
	calls struct {
		// ListTask holds details about calls to the ListTask method.
		ListTask []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
		}
	}
	lockListTask sync.RWMutex
}

// ListTask calls ListTaskFunc.
func (mock *ListTasksServiceMock) ListTask(ctx context.Context) (entity.Tasks, error) {
	if mock.ListTaskFunc == nil {
		panic("ListTasksServiceMock.ListTaskFunc: method is nil but ListTasksService.ListTask was just called")
	}
	callInfo := struct {
		Ctx context.Context
	}{
		Ctx: ctx,
	}
	mock.lockListTask.Lock()
	mock.calls.ListTask = append(mock.calls.ListTask, callInfo)
	mock.lockListTask.Unlock()
	return mock.ListTaskFunc(ctx)
}

// ListTaskCalls gets all the calls that were made to ListTask.
// Check the length with:
//
//	len(mockedListTasksService.ListTaskCalls())
func (mock *ListTasksServiceMock) ListTaskCalls() []struct {
	Ctx context.Context
} {
	var calls []struct {
		Ctx context.Context
	}
	mock.lockListTask.RLock()
	calls = mock.calls.ListTask
	mock.lockListTask.RUnlock()
	return calls
}

// Ensure, that AddTaskServiceMock does implement AddTaskService.
// If this is not the case, regenerate this file with moq.
var _ AddTaskService = &AddTaskServiceMock{}

// AddTaskServiceMock is a mock implementation of AddTaskService.
//
//	func TestSomethingThatUsesAddTaskService(t *testing.T) {
//
//		// make and configure a mocked AddTaskService
//		mockedAddTaskService := &AddTaskServiceMock{
//			AddTaskFunc: func(ctx context.Context, title string) (*entity.Task, error) {
//				panic("mock out the AddTask method")
//			},
//		}
//
//		// use mockedAddTaskService in code that requires AddTaskService
//		// and then make assertions.
//
//	}
type AddTaskServiceMock struct {
	// AddTaskFunc mocks the AddTask method.
	AddTaskFunc func(ctx context.Context, title string) (*entity.Task, error)

	// calls tracks calls to the methods.
	calls struct {
		// AddTask holds details about calls to the AddTask method.
		AddTask []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Title is the title argument value.
			Title string
		}
	}
	lockAddTask sync.RWMutex
}

// AddTask calls AddTaskFunc.
func (mock *AddTaskServiceMock) AddTask(ctx context.Context, title string) (*entity.Task, error) {
	if mock.AddTaskFunc == nil {
		panic("AddTaskServiceMock.AddTaskFunc: method is nil but AddTaskService.AddTask was just called")
	}
	callInfo := struct {
		Ctx   context.Context
		Title string
	}{
		Ctx:   ctx,
		Title: title,
	}
	mock.lockAddTask.Lock()
	mock.calls.AddTask = append(mock.calls.AddTask, callInfo)
	mock.lockAddTask.Unlock()
	return mock.AddTaskFunc(ctx, title)
}

// AddTaskCalls gets all the calls that were made to AddTask.
// Check the length with:
//
//	len(mockedAddTaskService.AddTaskCalls())
func (mock *AddTaskServiceMock) AddTaskCalls() []struct {
	Ctx   context.Context
	Title string
} {
	var calls []struct {
		Ctx   context.Context
		Title string
	}
	mock.lockAddTask.RLock()
	calls = mock.calls.AddTask
	mock.lockAddTask.RUnlock()
	return calls
}
