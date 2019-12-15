<template>
  <el-dialog :visible.sync="visible" :close-on-click-modal="false">
    <el-form>
      <el-form-item v-if="show.key" label="key">
        <el-input v-model="data.key" />
      </el-form-item>
      <el-form-item v-if="show.score" label="score">
        <el-input v-model="data.score" />
      </el-form-item>
      <el-form-item v-if="show.value" label="value">
        <el-button @click="jsonFormat">JSON格式化</el-button>
        <el-button @click="jsonCompress">JSON压缩</el-button>
        <el-input v-model="data.value" type="textarea" :autosize="{ minRows: 10 }" />
      </el-form-item>
    </el-form>

    <div slot="footer">
      <el-button type="primary">更新</el-button>
      <el-button @click="close">关闭</el-button>
    </div>
  </el-dialog>
</template>

<script>
export default {
  name: 'ValueEditDialog',
  data() {
    return {
      visible: false,
      data: {
        score: 0,
        key: '',
        value: ''
      },
      show: {
        score: false,
        key: false,
        value: false
      }
    }
  },
  methods: {
    open(data) {
      this.data = Object.assign({}, data)
      this.show = {
        score: !!data.score,
        key: !!data.key,
        value: !!data.value
      }
      this.visible = true
    },
    close() {
      this.data = {}
      this.show = {}
      this.visible = false
    },
    jsonFormat() {
      const json = JSON.parse(this.data.value)
      this.data.value = JSON.stringify(json, null, 2)
    },
    jsonCompress() {
      const json = JSON.parse(this.data.value)
      this.data.value = JSON.stringify(json)
    }
  }
}
</script>

<style scoped>

</style>
