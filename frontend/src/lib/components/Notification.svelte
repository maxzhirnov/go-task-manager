<script>
    import { successMessage } from '$lib/stores';
</script>

{#if $successMessage}
    <div class="notification success" role="alert">
        {$successMessage}
    </div>
{/if}

<style>
    .notification {
        position: fixed;
        top: 20px;
        right: 20px;
        padding: 1rem 1.5rem;
        border-radius: 4px;
        color: white;
        z-index: 1000;
        animation: slideIn 0.3s ease-out;
    }

    .success {
        background-color: #4CAF50;
        box-shadow: 0 2px 4px rgba(76, 175, 80, 0.3);
    }

    @keyframes slideIn {
        from {
            transform: translateX(100%);
            opacity: 0;
        }
        to {
            transform: translateX(0);
            opacity: 1;
        }
    }
</style>