import { onMounted, onUnmounted } from 'vue'

export function useScrollReveal(options = {}) {
  const {
    threshold = 0.1,
    rootMargin = '0px',
    className = 'is-visible'
  } = options

  let observer = null

  const initObserver = () => {
    observer = new IntersectionObserver((entries) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting) {
          entry.target.classList.add(className)
          
          // 如果只需要触发一次，观察后取消观察
          if (options.once) {
            observer.unobserve(entry.target)
          }
        } else if (!options.once) {
          entry.target.classList.remove(className)
        }
      })
    }, {
      threshold,
      rootMargin
    })

    // 观察所有带有 scroll-reveal 类的元素
    const elements = document.querySelectorAll('.scroll-reveal, .scroll-reveal-left, .scroll-reveal-right')
    elements.forEach((el) => observer.observe(el))
  }

  const cleanup = () => {
    if (observer) {
      observer.disconnect()
      observer = null
    }
  }

  onMounted(() => {
    initObserver()
  })

  onUnmounted(() => {
    cleanup()
  })

  return {
    initObserver,
    cleanup
  }
}

// 交错动画
export function useStaggerAnimation(selector = '.stagger-item', delay = 50) {
  onMounted(() => {
    const items = document.querySelectorAll(selector)
    items.forEach((item, index) => {
      item.style.animationDelay = `${index * delay}ms`
    })
  })
}

// 数字滚动动画
export function useCountUp(target, duration = 2000) {
  const animate = (element, start, end, duration) => {
    const startTime = performance.now()
    const range = end - start

    const step = (currentTime) => {
      const elapsed = currentTime - startTime
      const progress = Math.min(elapsed / duration, 1)
      
      // 使用缓动函数
      const easeOutQuart = 1 - Math.pow(1 - progress, 4)
      const current = start + range * easeOutQuart

      element.textContent = Math.floor(current)

      if (progress < 1) {
        requestAnimationFrame(step)
      } else {
        element.textContent = end
      }
    }

    requestAnimationFrame(step)
  }

  onMounted(() => {
    const elements = document.querySelectorAll(target)
    
    const observer = new IntersectionObserver((entries) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting) {
          const element = entry.target
          const endValue = parseInt(element.getAttribute('data-count') || element.textContent)
          animate(element, 0, endValue, duration)
          observer.unobserve(element)
        }
      })
    })

    elements.forEach((el) => observer.observe(el))
  })
}
