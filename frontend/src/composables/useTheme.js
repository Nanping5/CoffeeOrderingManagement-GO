import { ref, watch } from 'vue'

const THEME_KEY = 'coffee-theme'
const isDark = ref(false)
let initialized = false

// 应用主题到 DOM
const applyTheme = () => {
  const html = document.documentElement
  
  if (isDark.value) {
    html.setAttribute('data-theme', 'dark')
    html.classList.add('dark')
  } else {
    html.setAttribute('data-theme', 'light')
    html.classList.remove('dark')
  }
  
  localStorage.setItem(THEME_KEY, isDark.value ? 'dark' : 'light')
}

// 初始化主题（可在任何地方调用）
const initTheme = () => {
  if (initialized) return
  
  const savedTheme = localStorage.getItem(THEME_KEY)
  const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
  
  isDark.value = savedTheme ? savedTheme === 'dark' : prefersDark
  applyTheme()
  
  // 监听系统主题变化
  const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
  mediaQuery.addEventListener('change', (e) => {
    if (!localStorage.getItem(THEME_KEY)) {
      isDark.value = e.matches
      applyTheme()
    }
  })
  
  initialized = true
}

// 切换主题
const toggleTheme = () => {
  isDark.value = !isDark.value
  applyTheme()
}

// 设置主题
const setTheme = (theme) => {
  isDark.value = theme === 'dark'
  applyTheme()
}

// 监听 isDark 变化
watch(isDark, () => {
  applyTheme()
})

export function useTheme() {
  // 确保初始化
  if (!initialized && typeof window !== 'undefined') {
    initTheme()
  }
  
  return {
    isDark,
    toggleTheme,
    setTheme,
    initTheme
  }
}
