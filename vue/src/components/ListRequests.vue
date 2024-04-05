<template>
  <v-container class="fill-height">
    <v-responsive
      class="align-centerfill-height mx-auto"
      max-width="900"
    >

      <div class="text-center">
        <div class="text-body-2 font-weight-light mb-n1">Welcome to</div>

        <h1 class="text-h2 font-weight-bold">Listener!</h1>
      </div>

      <div class="text-center mt-3">
        <v-btn elevation="20" class="mr-5" @click="clear()">
          Clear
        </v-btn>
      </div>

      <div class="py-4" />

      <v-row>
        <v-col cols="12" v-if="loading">
          <div class="loading">Loading...</div>
        </v-col>

        <v-col cols="5" v-if="list">
          <v-list lines="two">
            <v-list-item
              v-for="item in list"
              :key="item.id"
              :title="'Request: ' + item.id"
              :subtitle="'Method ' + item.method + ' (date: ' + item.date + ')' "
            >
              <v-btn elevation="4" @click="() => loadElement(item.id)">
                show
              </v-btn>
            </v-list-item>
          </v-list>
        </v-col>
        <v-col cols="7">
          <div v-if="!selectedElement">Select element</div>
          <div v-if="loadingElement">Loading...</div>
          <div v-if="selectedElement">
            <v-card :title="selectedElement.uri"
                    :subtitle="'Method ' + selectedElement.method + ' (date: ' + selectedElement.created + ')'"
                    :text="'host: ' + selectedElement.host">
            </v-card>
            <v-card title="BODY" :text="selectedElement.body"
            ></v-card>
            <v-card title="HEADER">
              <v-list>
                <v-list-item v-for="(h, i) in selectedElement.header">
                  {{ i }}: {{ h }}
                </v-list-item>
              </v-list>
            </v-card>
          </div>
        </v-col>
      </v-row>
    </v-responsive>
  </v-container>
</template>

<script setup>
import { ref } from 'vue'
import axios from "axios";

const loading = ref(true)
const loadingElement = ref(false)
const list = ref([])
const selectedElement = ref(null)

async function clear() {
  try {
    await axios.get('http://localhost/remove');
    list.value = []
    selectedElement.value = null
  } catch (error) {
    console.error(error);
  }
}

async function getList() {
  try {
    let requestList = await axios.get('http://localhost/list');
    list.value = requestList.data
  } catch (error) {
    console.error(error);
  } finally {
    loading.value = false
  }
}

async function getElement(id) {
  try {
    let element = await axios.get('http://localhost/read/' + id);
    console.log(element.data[0])
    selectedElement.value = element.data[0]
  } catch (error) {
    console.error(error);
  } finally {
    loadingElement.value = false
  }
}

getList()

function loadElement(id) {
  loadingElement.value = true
  getElement(id)
}



</script>
