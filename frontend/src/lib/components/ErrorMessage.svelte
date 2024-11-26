<script>
    import { errorMessage } from '../stores.js';
</script>

{#if $errorMessage}
    <div class="error-message" role="alert">
        {$errorMessage}
    </div>
{/if}

<style>
    .error-message {
        color: red;
        text-align: center;
        margin-bottom: 10px;
        padding: 10px;
        background-color: #ffebee;
        border-radius: 4px;
        animation: fadeOut 3s forwards;
    }

    @keyframes fadeOut {
        0% { opacity: 1; }
        70% { opacity: 1; }
        100% { opacity: 0; }
    }
</style>