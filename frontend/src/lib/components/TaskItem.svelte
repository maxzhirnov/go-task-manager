<script>
    import { api } from '../api.js';
    import { tasks } from '../stores.js';
    import { showError } from '$lib/stores.js';


    export let task;
    let isEditing = false;
    let editTitle = task.title;
    let editDescription = task.description;
    let editStatus = task.status;

    async function handleStatusChange(newStatus) {
        try {
            await api.updateTaskStatus(task.id, newStatus);
            const updatedTasks = await api.fetchTasks();
            tasks.set(updatedTasks);
        } catch (error) {
            showError(error.message);
        }
    }

    async function handleDelete() {
        try {
            await api.deleteTask(task.id);
            const updatedTasks = await api.fetchTasks();
            tasks.set(updatedTasks);
        } catch (error) {
            showError(error.message);
        }
    }

    async function handleEdit() {
        if (editTitle.length > 100) {
            showError("Task title cannot exceed 100 characters.");
            return;
        }

        try {
            await api.updateTask(task.id, editTitle, editDescription, editStatus);
            const updatedTasks = await api.fetchTasks();
            tasks.set(updatedTasks);
            isEditing = false;
        } catch (error) {
            showError(error.message);
        }
    }

    function cancelEdit() {
        editTitle = task.title;
        editDescription = task.description;
        editStatus = task.status;
        isEditing = false;
    }
</script>
<li class="task-item {task.status}">
    {#if isEditing}
        <div class="task-edit">
            <input 
                bind:value={editTitle} 
                type="text" 
                placeholder="Task Title" 
                required
                maxlength="100"
            />
            <span class="character-count">({editTitle.length}/100)</span>
            <textarea 
                bind:value={editDescription} 
                placeholder="Task Description" 
                required
            ></textarea>
            <select bind:value={editStatus}>
                <option value="pending">Pending</option>
                <option value="in_progress">In Progress</option>
                <option value="completed">Completed</option>
            </select>
            <div class="edit-actions">
                <button class="save-btn" on:click={handleEdit}>Save</button>
                <button class="cancel-btn" on:click={cancelEdit}>Cancel</button>
            </div>
        </div>
    {:else}
        <div class="task-content">
            <h3>{task.title}</h3>
            <p>{task.description}</p>
        </div>
        <div class="task-actions">
            <div class="edit-actions">
                <button class="edit-btn" on:click={() => isEditing = true}>
                    Edit
                </button>
                <button class="delete-btn" on:click={handleDelete}>
                    Delete
                </button>
            </div>
            <div class="edit-status">
                <select 
                    value={task.status} 
                    on:change={(e) => handleStatusChange(e.target.value)}
                >
                    <option value="pending">Pending</option>
                    <option value="in_progress">In Progress</option>
                    <option value="completed">Completed</option>
                </select>
            </div>
        </div>
    {/if}
</li>

<style>
    .task-item {
        background-color: #f5f5f5;
        margin: 10px 0;
        padding: .7rem .7rem .7rem 1rem;
        border-radius: 5px;
        display: flex;
        justify-content: space-between;
        align-items: flex-start;
    }

    .task-content h3 {
        margin: 0 0 .7rem 0;
    }

    .task-content p {
        margin: 0 0 .7rem 0;
    }

    .edit-actions {
        display: flex;
        gap: 10px;
    }

    .edit-status { 
        margin-top: 1rem;
    }

    .edit-status select {
        width: 100%;
    }

    .task-edit {
        display: flex;
        flex-direction: column;
        gap: 10px;
        width: 100%;
    }

    .task-edit input,
    .task-edit textarea,
    .task-edit select {
        padding: 8px;
        border: 1px solid #ddd;
        border-radius: 4px;
    }

    .task-edit textarea {
        min-height: 100px;
        resize: vertical;
    }

    .edit-actions {
        display: flex;
        gap: 10px;
    }

    .character-count {
        font-size: 12px;
        color: #666;
        margin-left: 5px;
    }

    .save-btn,
    .cancel-btn,
    .edit-btn,
    .delete-btn {
        padding: 8px 16px;
        border: none;
        border-radius: 4px;
        cursor: pointer;
        font-size: 14px;
    }

    .save-btn {
        background-color: #4CAF50;
        color: white;
    }

    .cancel-btn {
        background-color: #757575;
        color: white;
    }

    .edit-btn {
        background-color: #2196F3;
        color: white;
    }

    .delete-btn {
        background-color: #f44336;
        color: white;
    }

    .task-item.pending { border-left: 6px solid #ffd700; }
    .task-item.in_progress { border-left: 6px solid #1e90ff; }
    .task-item.completed { border-left: 6px solid #32cd32; }
</style>