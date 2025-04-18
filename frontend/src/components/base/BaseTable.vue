<template>
  <div class="table-wrapper">
    <slot></slot>
  </div>
</template>

<style scoped>
.table-wrapper {
  width: 100%;
  background: transparent;
  border-radius: 8px;
  border: 2px solid #646464;
  box-shadow: 0 4px 6px rgba(100, 100, 100, 0.2), 0 8px 12px rgba(100, 100, 100, 0.1);
  /* Eliminar el margin-bottom que puede estar causando espacio adicional */
  margin-bottom: 0;
  overflow: hidden;
}
</style>