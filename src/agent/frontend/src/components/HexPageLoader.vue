<template>
    <transition name="fade">
        <div class="ploader-main" transition="fade-transition" v-if="!isLoaded">
            <div class="ploader-container">
                <div></div>
                <div></div>
                <div></div>
                <div></div>
                <div></div>
                <div></div>
                <div></div>
                <div></div>
                <div></div>
                <div></div>
                <div></div>
                <div></div>
                <div></div>
                <div></div>
            </div>
        </div>
    </transition>
</template>

<script>
export default {
    data: () => ({
        isLoaded: false,
    }),
    
    mounted() {
        document.onreadystatechange = () => {
            if(document.readyState == "complete") {
                this.isLoaded = true;
            }
        };
    },
}
</script>

<style lang="scss" scoped>
.fade-enter-active, .fade-leave-active {
    transition: opacity .5s;
}

.fade-enter, .fade-leave-to {
    opacity: 0;
}

.ploader-main {
    display: grid;
    grid-template-columns: repeat(auto-fit,calc(var(--s) + 2*var(--mh)));
    grid-template-areas: 
    "hexd hexd hexd hexd hexd"
    "hexd hexd hexd hexd hexd"
    "hexd hexd hexd hexd hexd"
    "hexd hexd hexd hexd hexd"
    "hexd hexd hexd hexd hexd"
    "hexd hexd hexd hexd hexd"
    "hexd hexd hexd hexd hexd";
    justify-content:center;
    padding-top: 40vh;
    background: var(--v-secondary-base);
    height: 100vh;
    width: 100vw;

    --s: 80px; /* size */
    --r: 1.15; /* ratio */
    /* clip-path */
    --h: 0.5;  
    --v: 0.25; 
    --hc:calc(clamp(0,var(--h),0.5) * var(--s)) ;
    --vc:calc(clamp(0,var(--v),0.5) * var(--s) * var(--r)); 
    
    /*margin */
    --mv: 2px; /* vertical */
    --mh: calc(var(--mv) + (var(--s) - 2*var(--hc))/2); /* horizontal */
    /* for the float*/
    --f: calc(2*var(--s)*var(--r) + 4*var(--mv)  - 2*var(--vc) - 2px);
  }
  
  .ploader-container {
    grid-column: 1/-1;
    max-width:790px;
    margin:0 auto;
    font-size: 0; /*disable white space between inline block element */
    position:relative;
    padding-bottom:50px;
    // filter:drop-shadow(2px 2px 1px #333)
    grid-area: hexd;
  }
  
  .ploader-container div {
    width: var(--s);
    margin: var(--mv) var(--mh);
    height: calc(var(--s)*var(--r)); 
    display: inline-block;
    font-size:initial;
    clip-path: polygon(var(--hc) 0, calc(100% - var(--hc)) 0,100% var(--vc),100% calc(100% - var(--vc)), calc(100% - var(--hc)) 100%,var(--hc) 100%,0 calc(100% - var(--vc)),0 var(--vc));
    margin-bottom: calc(var(--mv) - var(--vc)); 
  }
  
  .ploader-container::before{
    content: "";
    width: calc(var(--s)/2 + var(--mh));
    float: left;
    height: 120%;
    shape-outside: repeating-linear-gradient(     
                     transparent 0 calc(var(--f) - 2px),      
                     #fff        0 var(--f));
  }
  
  .ploader-container div::before {
    padding-top:20px;
    content:"Salt\A Agent";
    text-transform:uppercase;
    white-space:pre;
    font-size:70px;
    font-family:'Poppins' sans-serif;
    font-weight:400;
    text-align:center;
    position:absolute;
    color:#e9e9ef;
    // background:linear-gradient(45deg,#55636e,#d28994);
    background: radial-gradient(circle, var(--v-accent-lighten2), var(--v-primary-base));
    inset:0;
  }
  
  .ploader-container div {
    animation:show 1.5s infinite;
    opacity:0;
  }
  @for $i from 1 through 27 {
      .ploader-container div:nth-child(#{random(28)}) {
       animation-delay:(1.5*random())*1s
    }
  }
  
  @keyframes show{
    60% {
       opacity:1;
    }
  }
  
  body  {
    background:#e9e9ef;
  }
</style>