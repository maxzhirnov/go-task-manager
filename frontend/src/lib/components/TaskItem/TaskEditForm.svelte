<!-- TaskEditForm.svelte -->
<script>
    import StatusSelector from './StatusSelector.svelte';
    
    // Props
    export let editTitle = '';
    export let editDescription = '';
    export let editStatus;
    export let onSave;
    export let onCancel;
    export let statusSelector;

    function handleSave() {
        if (onSave) onSave({
            title: editTitle,
            description: editDescription,
            status: editStatus
        });
    }
</script>

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
    
    <StatusSelector 
        bind:this={statusSelector}
        value={editStatus}
        onChange={(newStatus) => editStatus = newStatus}
        immediate={false}
        fullWidth={true}
    />
    
    <div class="edit-actions-in-edit">
        <button class="tech-button save" on:click={handleSave}>
            <span class="button-content">
                <span class="icon">⚡</span>
                <span class="text">SAVE_CHANGES</span>
            </span>
            <div class="button-glow"></div>
        </button>
        
        <button class="tech-button cancel" on:click={onCancel}>
            <span class="button-content">
                <span class="icon">×</span>
                <span class="text">ABORT</span>
            </span>
            <div class="button-glow"></div>
        </button>
    </div>
</div>

<style>
    .task-edit {
        width: 100%;
        margin-top: 1rem;
        display: flex;
        flex-direction: column;
        gap: 0.8rem;
    }

    .task-edit input,
    .task-edit textarea {
        background: #2d3436;
        border: 1px solid #0984e3;
        color: #fff;
        padding: 0.5rem;
        border-radius: 3px;
        font-family: "JetBrains Mono", monospace;
    }
    
    button {
        background: transparent;
        border: 1px solid #00b894;
        color: #00b894;
        padding: 0.4rem;
        border-radius: 3px;
        cursor: pointer;
        transition: all 0.3s ease;
    }

    button:hover {
        background: rgba(0, 184, 148, 0.1);
        box-shadow: 0 0 8px rgba(0, 184, 148, 0.3);
    }

    .save {
        color: #00b894;
        border-color: #00b894;
    }

    .save:hover {
        background: rgba(0, 184, 148, 0.1);
        box-shadow: 0 0 10px rgba(0, 184, 148, 0.3);
    }

    .save::after {
        content: '';
        position: absolute;
        top: -2px;
        left: -2px;
        right: -2px;
        bottom: -2px;
        background: linear-gradient(45deg, #00b894, transparent);
        z-index: -1;
        opacity: 0;
        transition: opacity 0.3s ease;
    }

    /* Cancel button styling */
    .cancel {
        color: #ff6b6b;
        border-color: #ff6b6b;
    }

    .cancel:hover {
        background: rgba(255, 107, 107, 0.1);
        box-shadow: 0 0 10px rgba(255, 107, 107, 0.3);
    }

    .cancel::after {
        content: '';
        position: absolute;
        top: -2px;
        left: -2px;
        right: -2px;
        bottom: -2px;
        background: linear-gradient(45deg, #ff6b6b, transparent);
        z-index: -1;
        opacity: 0;
        transition: opacity 0.3s ease;
    }
</style>