<template>
  <el-row class="app-container">
    <el-col :span="8">
      <el-select>
        <el-option>Redis 1</el-option>
        <el-option>Redis 2</el-option>
      </el-select>
      <el-button icon="el-icon-circle-plus-outline"></el-button>
      <el-input/>
      <el-tree ref="tree" :data="keys" class="filter-tree"/>
    </el-col>

    <el-col :span="16">
      <el-form :data="data" label-width="80px">
        <el-form-item :label="data.type">
          <el-input v-model="data.key"/>
        </el-form-item>
      </el-form>
    </el-col>
  </el-row>
</template>


<script>
export default {
  data() {
    return {
      keys: [],
      data: {
        type: 'STRING',
        key: '123:456',
        value: 'asdf'
      }
    }
  },
  methods: {
  }
}
</script>

<style scoped>
.menu {
  width: 300px;
  height: 300px;
  background: yellowgreen;
  display: inline;
}
.container {
  display: inline;
}
</style>


