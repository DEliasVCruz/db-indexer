<script setup lang="ts">
import TheFooter from "@/components/TheFooter.vue";
import TheHeader from "@/components/TheHeader.vue";
import AboutSection from "@/components/AboutSection.vue";
</script>

<template>
  <div class="outer-layer">
    <TheHeader :side-bar="false" :search-bar="false" />
    <main class="my-1 flex flex-col gap-5 overflow-auto">
      <h1 class="mx-5 align-middle text-4xl font-bold sm:ml-[22%]">About</h1>
      <AboutSection title="The Usage">
        <p>
          This app is meant to be used as simple database indexer for mailing
          directories. A <b>mailing directory</b> being a set of folders
          containing mail text files.<br />
          This files contain two important sets of data for this application,
          the <i>metadata</i> and the <b>contents</b>. Though the contents of
          the mail is what most users will usually be interested on, it's
          important that the file contains the proper metadata to be properly
          indexed.<br />
          For the indexing of a given mail to be completed it is required that
          the <code>x_filename</code> metadata tag is present, otherwise the
          indexing of the given mail will be aborted and continue with the next
          email.<br />
          You can upload your mailing directories from the landing page using a
          compressed format like gzip, zip or tar with a maximum file size of
          500 megabytes.<br />
          Once the indexing is complete you can <b>full text search</b> the
          contents of your indexed mails as well as the contents of specific
          fields through the search bar.
        </p>
      </AboutSection>
      <AboutSection title="The Developer">
        <p>
          My name is Daniel Vilela and I'm a 24 year's old aspiring software
          engineer. Though I graduated in business administration, later on my
          career I realized that computer science and programming are my
          passion.<br />
          I'm <b>self thought</b> on all things related to CS and programming,
          through relentless dedication and commitment I managed to teach my
          self the fundamentals of computers systems and programming
          languages.<br />
          Recently I managed to land my first tech related job as a data
          engineer in a renown company in my home country, by showcasing my
          skills as a programmer and my self thought nature and approach to
          problem solving. Though it has been an incredible opportunity to get
          my feet wet with the ins and outs of a tech related position and the
          data engineering position served as a nice breach between my business
          background and my programming skills, I now feel ready to start amino
          for more software engineer related jobs<br />
          Though I'm most proficient with <b>Python</b> and
          <b>Linux core utils</b>, I have recently taken a linking for
          <b>Go</b> and web applications. I am a heavy <b>neovim</b> user and is
          my editor of choice (you can find my configuration on
          <a href="https://github.com/DEliasVCruz/CristalMoon">Cristal Moon</a
          >). I'm also a Linux user and enthusiast and up until recently
          <i>I use Arch btw</i> but due to job requirements I'm currently
          working on WSL (which is not that bad and makes the Windows experience
          manageable). You can take a look at my
          <a href="https://github.com/DEliasVCruz/dotfiles">dotfiles</a> if you
          want.
        </p>
      </AboutSection>
      <AboutSection title="The Motivation">
        <p>
          As mentioned earlier in this page, I have recently started to learn go
          and a job application required me to build a tool like this. So I took
          the opportunity to make this a side project and turn it into a full
          blown application.<br />
          I feel very proud about my accomplishment (for someone who had never
          learn about Go, Vue and Tailwind until the beginning of this project)
          and have decided to make it a part of my portfolio.<br />
          As such, I will continue to develop this application to a point i feel
          satisfied with it's state and share it to the world.<br />
          If you've this application interesting or helpful and would like to
          know more about the author you can follow me on
          <a href="https://www.linkedin.com/in/daniel-vilela-cruz-694968138/"
            >Linkedin</a
          >
          or my <a href="https://github.com/DEliasVCruz">Github</a> profile
        </p>
      </AboutSection>
    </main>
    <TheFooter />
  </div>
</template>

<style scoped></style>
