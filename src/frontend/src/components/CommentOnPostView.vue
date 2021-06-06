<template>
  <div class="post-comment">
    <div class="post-comment-left-side">
      <v-img  class="post-comment-profile"
              src="https://i.pinimg.com/736x/4d/8e/cc/4d8ecc6967b4a3d475be5c4d881c4d9c.jpg"
              @click="$router.push('/user')"
              alt="Profile picture"/>
      <b @click="$router.push('/user')" style="cursor: pointer">Username1</b>
      <p style="margin: 0 0 0 5px; font-size: 13px; text-align: left;">My First comment askj dsajdaskjds a;kj d;kaslk ja dksadkjasdk;jsakjd askjdhsakjh daskjhdsa;jhkdsa;kjhd djsahjhldkaskjlhdas kjhlasjkhldasjklhdaskjlhdaslkjhdlkjhasdlkjh asdkjlhaskjhdas;kjh ;dkjh as;jkhdsaj;khdj;khasd;kjhsa ;kjhdsa;kjhd;kjhsad;kjhasdjhka jashkjh;kdsaas;jhkdasj;hkd;jahskd;jhasasjh;d j;hasjh; sajh;d jhsajhd asjh;d asjh;jh; sa;jhdasj;hd sa;jhsdaskjdaskjhdksl ajdkljaskdlj sakjld;askl;das;k lhdask;lhd;klahsjd;klhjsakdh;ljalsd jklasdkjaskj adjkb;kjsad;j klsad;as;jds a;jbdas;j dsj;b</p>
      <p style="text-align: left; font-size: 10pt; margin: 0 0 0 5px; color: #858585; width: 120px">
        Posted 1 hour ago
      </p>
    </div>

  </div>
</template>

<script>
export default {
  name: "CommentOnPostView",

}
</script>

<style scoped>

.post-comment {
  min-height: 10px;
  height: auto;
  max-height: 50px;
}

.post-comment-left-side {
  direction: ltr;

  flex-direction: row;
  text-align: -webkit-center;
  align-items: center;
  float: left;
  display: flex;
  justify-content: left;
  overflow-x: hidden;
  overflow-y: hidden;

}

.post-comment-profile {
  width: 30px;
  height: 30px;
  object-fit: cover;
  border-radius: 20%;
  margin: 3px;
  cursor: pointer;


  filter: brightness(1);

  transition: .3s;
  z-index: 0;
}

</style>