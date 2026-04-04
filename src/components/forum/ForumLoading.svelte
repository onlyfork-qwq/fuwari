<script lang="ts">
export let text = "正在加载...";
export let fullHeight = true;
</script>

<div
	class="forum-loading flex flex-col items-center justify-center gap-4 {fullHeight
		? 'min-h-[60vh]'
		: ''}"
>
	<div class="loading-spinner">
		<div class="spinner-ring"></div>
		<div class="spinner-ring"></div>
		<div class="spinner-ring"></div>
	</div>
	<p class="text-white/50 text-sm animate-pulse">{text}</p>
</div>

<style>
	.loading-spinner {
		position: relative;
		width: 48px;
		height: 48px;
	}

	.spinner-ring {
		position: absolute;
		width: 100%;
		height: 100%;
		border-radius: 50%;
		border: 3px solid transparent;
		animation: spin 1.2s cubic-bezier(0.5, 0, 0.5, 1) infinite;
	}

	.spinner-ring:nth-child(1) {
		border-top-color: var(--primary);
		animation-delay: -0.3s;
	}

	.spinner-ring:nth-child(2) {
		border-right-color: var(--primary);
		animation-delay: -0.15s;
		opacity: 0.7;
	}

	.spinner-ring:nth-child(3) {
		border-bottom-color: var(--primary);
		opacity: 0.4;
	}

	@keyframes spin {
		0% {
			transform: rotate(0deg);
		}
		100% {
			transform: rotate(360deg);
		}
	}

	.forum-loading {
		animation: fade-in 0.3s ease-out;
	}

	@keyframes fade-in {
		from {
			opacity: 0;
		}
		to {
			opacity: 1;
		}
	}
</style>