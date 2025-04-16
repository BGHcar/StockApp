<template>
  <div class="table-wrapper">
    <slot></slot>
  </div>
</template>

<style scoped>
.table-wrapper {
  width: 100%;
  background: rgba(60, 16, 83, 0.95);
  border-radius: 8px;
  border: 2px solid rgba(173, 83, 137, 0.3);
  box-shadow: 
    0 4px 6px rgba(60, 16, 83, 0.3),
    0 8px 12px rgba(173, 83, 137, 0.2),
    0 0 0 2px rgba(173, 83, 137, 0.1);
}
</style>