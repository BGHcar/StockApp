<template>
  <div class="base-scroll">
    <slot></slot>
  </div>
</template>

<style scoped>
.base-scroll {
  overflow-y: auto;
}

.base-scroll::-webkit-scrollbar {
  width: 10px;
}

.base-scroll::-webkit-scrollbar-track {
  background: rgba(60, 16, 83, 0.8);
  border-radius: 15px;
}

.base-scroll::-webkit-scrollbar-thumb {
  background: rgba(173, 83, 137, 0.3);
  border-radius: 15px;
}

.base-scroll::-webkit-scrollbar-thumb:hover {
  background: rgba(173, 83, 137, 0.5);
}
</style>