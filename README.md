# Khamsa Application's Backend

## Overview

This is a Go-based backend application for the Khamsa App built for Salam Hackathon.

## Architecture

The application follows a RESTful API architecture built with the Gin web framework. It connects to Firestore for data persistence and Google's Generative AI for content generation.

### Core Components

1. **Server**: The main server component that handles HTTP requests
2. **Firestore**: Database integration for storing user sessions, learning paths, projects, and tasks
3. **Generative AI**: AI integration for creating learning content and providing help

## Data Models

The application works with the following key data models:

### Learning

A learning path that contains multiple projects.

```go
type Learning struct {
    Id          string    // Unique identifier
    Language    string    // Programming language (e.g., "JavaScript")
    Level       string    // User's skill level (e.g., "Beginner")
    FrameWork   string    // Framework to learn (e.g., "React")
    Goal        string    // Learning objective
    Title       string    // Title of the learning path
    Description string    // Description of the learning path
    Progress    int       // Progress tracked as completed projects
}
```

### Project

A project within a learning path, containing multiple tasks.

```go
type Project struct {
    Id          string    // Unique identifier
    Order       int       // Order within the learning path
    Title       string    // Project title
    Description string    // Project description
    IsLocked    bool      // Whether the project is available to the user
    LearningId  string    // Reference to parent learning path
    Progress    int       // Number of completed tasks
    TaskNum     int       // Total number of tasks
}
```

### Task

A specific task within a project.

```go
type Task struct {
    Id          string    // Unique identifier
    Order       int       // Order within the project
    Title       string    // Task title
    Description string    // Task description
    Completed   bool      // Whether the task is completed
    ProjectId   string    // Reference to parent project
}
```

## API Endpoints

### Learning Paths

#### `POST /new-learning/:id`

Creates a new learning path for a user session.

**Path Parameters:**
- `id`: Session ID

**Request Body:**
```json
{
    "language": "JavaScript",
    "level": "Beginner",
    "framework": "React",
    "goal": "Build a personal portfolio website"
}
```

**Response:**
A complete learning path with projects and tasks.

#### `GET /learnings/:id`

Gets all learning paths for a user session.

**Path Parameters:**
- `id`: Session ID

**Response:**
Array of learning paths with progress information.

### Projects

#### `POST /projects/:id`

Gets all projects for a specific learning path.

**Path Parameters:**
- `id`: Session ID

**Request Body:**
```json
{
    "id": "learning-path-id"
}
```

**Response:**
Array of projects with progress information.

### Tasks

#### `POST /tasks/:id`

Gets all tasks for a specific project.

**Path Parameters:**
- `id`: Session ID

**Request Body:**
```json
{
    "id": "project-id"
}
```

**Response:**
Array of tasks.

#### `POST /tasks/check/:id`

Marks a task as completed and potentially unlocks the next project.

**Path Parameters:**
- `id`: Session ID

**Request Body:**
```json
{
    "id": "task-id"
}
```

**Response:**
Confirmation message and information about unlocked projects.

### Help and Suggestions

#### `POST /help`

Requests AI-generated help for a specific programming task.

**Request Body:**
```json
{
    "framework": "React",
    "language": "JavaScript",
    "task": "Implement form validation",
    "project": "User registration system"
}
```

**Response:**
AI-generated help content.

#### `POST /suggest`

Requests AI-generated learning path suggestions based on user preferences.

**Request Body:**
```json
{
    "preference": "Web Development",
    "level": "Intermediate",
    "goal": "Build a full-stack web application"
}
```

**Response:**
AI-generated suggestions for learning paths.

## Progression System

The application implements a progression system where:

1. Learning paths contain multiple projects
2. Projects contain multiple tasks
3. Tasks must be completed to make progress
4. Projects are initially locked except for the first one
5. Completing all tasks in a project unlocks the next project
6. Progress is tracked at both project and learning path levels
