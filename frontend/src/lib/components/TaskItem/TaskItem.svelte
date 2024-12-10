<script>
    import { api } from '../../api.js';
    import { tasks } from '../../stores.js';
    import { showError } from '$lib/stores.js';
    import { dragHandle } from 'svelte-dnd-action';
    import Time from "svelte-time";
    import { Analytics } from '$lib/analytics';

    import FormattedTime from './FormattedTime2.svelte';
    import TechButton from './TechButton.svelte';
    import DragHandle from './DragHandle.svelte';
    import StatusSelector from './StatusSelector.svelte';
    import TaskDescription from './TaskDescription.svelte';
    import TaskTitle from './TaskTitle.svelte';
    import TaskEditForm from './TaskEditForm.svelte';

    export let task;

    let isEditing = false;
    let editTitle = task.title;
    let editDescription = task.description;
    let editStatus = task.status;
    let statusSelector; 

    async function handleStatusChange(newStatus) {
        try {
            await api.updateTaskStatus(task.id, newStatus);
            const updatedTasks = await api.fetchTasks();
            tasks.set(updatedTasks);

            Analytics.track('Task Status Updated', {
                taskId: task.id,
                status: newStatus
            });

        } catch (error) {
            showError(error.message);
        }
    }

    async function handleDelete() {
        try {
            await api.deleteTask(task.id);
            const updatedTasks = await api.fetchTasks();
            tasks.set(updatedTasks);

            Analytics.track('Task Deleted', {
                taskId: task.id
                });
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
            const newStatus = statusSelector.getCurrentValue();
            await api.updateTask(task.id, editTitle, editDescription, newStatus);
            const updatedTasks = await api.fetchTasks();
            tasks.set(updatedTasks);
            isEditing = false;
            
            Analytics.track('Task Edited', {
                taskId: task.id,
                title: editTitle,
                description: editDescription,
                status: newStatus
            });
        } catch (error) {
            showError(error.message);
        }
    }

    function cancelEdit() {
        editTitle = task.title;
        editDescription = task.description;
        editStatus = task.status;
        isEditing = false;
        Analytics.track('Task Edit Cancelled', {
            taskId: task.id
        });
    }
</script>

<div class="task-item {task.status}">
    <div class="drag-handle" use:dragHandle><DragHandle /></div>
        {#if isEditing}
            <TaskEditForm
                bind:editTitle
                bind:editDescription
                bind:editStatus
                bind:statusSelector
                onSave={handleEdit}
                onCancel={() => isEditing = false}
            />
        {:else}
            <div class="task-content-container">
                <div class="task-content">
                    <div class="task-text">
                        <div class="task-text-content">
                            <TaskTitle title={task.title} />
                            {#if task.description}
                                <TaskDescription description={task.description} taskId={task.id} />
                            {/if}
                        </div>
                        <div class="task-text-footer">
                            <FormattedTime timestamp={task.created_at} daysToReRelative={7} />
                        </div>
                    </div>
                </div>
                <div class="task-actions">
                    <div class="edit-actions">
                        <TechButton 
                            type="edit" 
                            title="Edit task" 
                            on:click={() => isEditing = true} 
                        />
                        <TechButton 
                            type="delete" 
                            title="Delete task" 
                            on:click={handleDelete} 
                        />
                        <TechButton 
                            type="complete" 
                            title="Complete task" 
                            on:click={() => handleStatusChange("completed")} 
                        />
                    </div>
                    <div class="edit-status">
                        <div class="edit-status">
                            <StatusSelector 
                                value={task.status} 
                                onChange={handleStatusChange}
                                immediate={true}
                            />
                        </div>
                    </div>
                </div>
            </div>
        {/if}
    </div>

<style>
    .task-item {
        display: flex;
        width: 100%;
        background: #1c1c1c;
        border: 1px solid #0984e3;
        border-radius: 4px;
        padding: 0 1rem 1rem 0;
        margin-bottom: 1rem;
        position: relative;
        overflow: hidden;
        box-shadow: 0 2px 8px rgba(9, 132, 227, 0.1);
        font-family: "JetBrains Mono", monospace;
    }

    /* Circuit board pattern background */
    .task-item::before {
        content: "";
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background-image: 
            radial-gradient(
                circle at 50% 50%,
                rgba(0, 184, 148, 0.05) 1px,
                transparent 1px
            );
        background-size: 10px 10px;
        pointer-events: none;
    }

    /* Status indicators */
    .task-item.pending::after {
        content: "STATUS::PENDING";
        position: absolute;
        top: 0.5rem;
        right: 0.5rem;
        font-size: 0.6rem;
        color: #ffd32a;
        opacity: 0.8;
    }

    .task-item.in_progress::after {
        content: "STATUS::RUNNING";
        position: absolute;
        top: 0.5rem;
        right: 0.5rem;
        font-size: 0.6rem;
        color: #0984e3;
        opacity: 0.8;
    }

    .task-item.completed::after {
        content: "STATUS::COMPLETE";
        position: absolute;
        top: 0.5rem;
        right: 0.5rem;
        font-size: 0.6rem;
        color: #00b894;
        opacity: 0.8;
    }

    .task-content-container {
        width: 100%;
        position: relative;
        z-index: 1;
    }

    .task-content {
        display: flex;
        gap: .2rem;
        align-items: flex-start;
    }

    .task-text {
        width: 100%;
    }

    .task-text-content {
        margin-top: .5rem;
        display: flex;
        flex-direction: column;
        gap: 0.8rem;
    }

    .edit-actions-in-edit {
        display: flex;
        gap: 1rem;
        margin-top: 1rem;
    }

    .tech-button {
        position: relative;
        background: transparent;
        border: 1px solid currentColor;
        padding: 0.8rem 1.5rem;
        font-family: 'JetBrains Mono', monospace;
        font-size: 0.9rem;
        cursor: pointer;
        overflow: hidden;
        transition: all 0.3s ease;
        flex: 1;
    }

    .button-content {
        position: relative;
        z-index: 2;
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 0.5rem;
    }

    .icon {
        font-size: 1.1em;
    }

    .button-glow {
        position: absolute;
        top: 0;
        left: -100%;
        width: 100%;
        height: 100%;
        background: linear-gradient(
            90deg,
            transparent,
            rgba(255, 255, 255, 0.1),
            transparent
        );
        transition: 0.5s;
    }

    @keyframes scan {
        to {
            transform: translateY(100%);
        }
    }

    /* Reduced motion */
    @media (prefers-reduced-motion: reduce) {
        .button-glow,
        .tech-button::before {
            display: none;
        }
    }

    .task-text-footer {
        margin-top: 0.5rem;
        font-size: 0.7rem;
        color: #636e72;
    }

    .task-actions {
        display: flex;
        justify-content: space-between;
        margin-top: 1rem;
        padding-top: 1rem;
        border-top: 1px solid rgba(9, 132, 227, 0.2);
    }


    .character-count {
        color: #636e72;
        font-size: 0.7rem;
        text-align: right;
    }

    /* Hover effects */
    .task-item:hover {
        border-color: #00b894;
        box-shadow: 0 0 15px rgba(0, 184, 148, 0.1);
    }

    .task-item:hover::before {
        animation: circuitPulse 2s infinite;
    }

    @keyframes circuitPulse {
        0%, 100% {
            opacity: 0.05;
        }
        50% {
            opacity: 0.1;
        }
    }

    /* Status-specific glows */
    .task-item.pending {
        border-left: 3px solid #ffd32a;
    }

    .task-item.in_progress {
        border-left: 3px solid #0984e3;
    }

    .task-item.completed {
        border-left: 3px solid #00b894;
    }

    
</style>