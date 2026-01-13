import { reactive } from 'vue'

const STORAGE_KEY = 'collected_articles'

const getInitialCollections = () => {
    try {
        const stored = localStorage.getItem(STORAGE_KEY)
        return stored ? JSON.parse(stored) : {}
    } catch (e) {
        return {}
    }
}

export const collectionStore = reactive({
    // Use an object for reactivity: { sn: true }
    collectedMap: getInitialCollections(),

    add(sn) {
        this.collectedMap[String(sn)] = true
        this.persist()
    },

    remove(sn) {
        delete this.collectedMap[String(sn)]
        this.persist()
    },

    isCollected(sn) {
        return !!this.collectedMap[String(sn)]
    },

    setCollections(snList) {
        const newMap = {}
        snList.forEach(sn => {
            newMap[String(sn)] = true
        })
        this.collectedMap = newMap
        this.persist()
    },

    persist() {
        localStorage.setItem(STORAGE_KEY, JSON.stringify(this.collectedMap))
    }
})
