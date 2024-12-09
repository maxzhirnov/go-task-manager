<script>
    export let title;
</script>

<div class="task-title-content">
    <h3 class="task-title">
            <span class="title-prefix">TASK::</span>
        {title}
    </h3>
</div>

<style>
    .title-prefix {
        color: #0984e3;
        font-weight: normal;
        font-size: .6rem;
        opacity: 0.8;
    }

    .task-title {
        font-family: "JetBrains Mono", monospace;
        font-size: 1rem;
        color: #fff;
        margin: 0;
        padding: 0.5rem;
        /* background: rgba(9, 132, 227, 0.1); */
        border-radius: 3px;
        position: relative;
        overflow: hidden;
    }

    .task-title:hover {
        background: rgba(9, 132, 227, 0.03);
    }

    .task-title:hover::before {
        opacity: 1;
    }


    /* Glitch effect on hover */
    .task-title:hover {
        animation: glitch 0.15s ease;
    }

    @keyframes glitch {
        0% { transform: translate(0); }
        20% { transform: translate(-1px, 1px); }
        40% { transform: translate(-1px, -1px); }
        60% { transform: translate(1px, 1px); }
        80% { transform: translate(1px, -1px); }
        100% { transform: translate(0); }
    }

    @media screen and (max-width: 400px) {
        .task-title {
            font-size: 0.8rem;
        }
    }
</style>