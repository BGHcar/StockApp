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
  /* Cambiado del color morado al gris que complementa el gradiente */
  background: rgba(100, 100, 100, 0.2);
  border-radius: 15px;
}

.base-scroll::-webkit-scrollbar-thumb {
  /* Cambiado del color rosa al gris más oscuro */
  background: rgba(100, 100, 100, 0.5);
  border-radius: 15px;
}

.base-scroll::-webkit-scrollbar-thumb:hover {
  /* Color de hover más oscuro */
  background: rgba(100, 100, 100, 0.8);
}
</style>