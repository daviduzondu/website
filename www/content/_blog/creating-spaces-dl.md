---
title: Creating a Twitter Spaces Downloader Using Node.js and FFmpeg
description: 
draft: false
tags: [javascript, webdev]
date: 28 June 2024
---

![virgin api consumer vs chad independent scraper](https://i.imgur.com/j4gD4fw.jpg)
A meme I found on 4chan 

Hello, My name is David Uzondu and I am a final year Computer Science student at Bayero University Kano. I was recently approved for the back-end track on the HNG internship program.

One reason I am excited to participate in the [HNG internship program](https://hng.tech/internship) is to gain hands-on experience and deepen my skills in back-end development. This opportunity will work on real-world projects, learn best practices, and collaborate with industry professionals. I also signed up for [HNG Premium](https://hng.tech/premium) which will give me access to exclusive opportunities like certificates, reference letters, and the latest job openings.

A few weeks ago, I wanted to listen to a Twitter Space conversation but had network issues. I noticed the space was recorded, so I bookmarked the link to return when my connection was stable. 

The next day, I tried to listen again, but the Twitter app wouldn’t load the space, which frustrated me. Since Twitter doesn’t allow downloading recorded spaces for later, I decided to use `yt-dlp`—a command-line tool that downloads media from sites like YouTube, Vimeo, and Twitter. However, there was a bug in the version of `yt-dlp` installed on my machine that caused the program to crash whenever I tried to download a space.

So I figured, why not create my own space downloader tool? It would be great practice for Node.js and TypeScript.

## Reverse-engineering the Twitter Internal API
On my Chrome desktop, I loaded the Twitter Space page while listening for the network requests in the DevTools. When you open a recorded Twitter Space link in Chrome or any other browser, here's what happens (assuming you are signed in):

Twitter loads a modal on the `/peek` route. This modal shows important information about the space, like the host and speakers, and provides a big blue "Play recording" button.

When the button is clicked, a GET request is sent to `https://x.com/i/api/graphql/SL4eyLXdr1zWZVpXRhxZ4Q?variables=<variables>&features=<features>`, where `variables` and `features` are query string parameters. `variables` is encoded with the following key-value pair:

```json
{
  "id": "<valid_twitter_spaces_id>",
  "isMetatagsQuery": true,
  "withReplays": true,
  "withListeners": true
}
```

And `features` is encoded with the following:
```json
{
  "rweb_tipjar_consumption_enabled": true,
  "responsive_web_graphql_exclude_directive_enabled": true,
  "verified_phone_label_enabled": false,
  "creator_subscriptions_tweet_preview_api_enabled": true,
  "responsive_web_graphql_timeline_navigation_enabled": true,
  "responsive_web_graphql_skip_user_profile_image_extensions_enabled": false,
  "communities_web_enable_tweet_community_results_fetch": true,
  "c9s_tweet_anatomy_moderator_badge_enabled": true,
  "articles_preview_enabled": true,
  "tweetypie_unmention_optimization_enabled": true,
  "responsive_web_edit_tweet_api_enabled": true,
  "graphql_is_translatable_rweb_tweet_is_translatable_enabled": true,
  "view_counts_everywhere_api_enabled": true,
  "longform_notetweets_consumption_enabled": true,
  "responsive_web_twitter_article_tweet_consumption_enabled": true,
  "tweet_awards_web_tipping_enabled": false,
  "creator_subscriptions_quote_tweet_preview_enabled": false,
  "freedom_of_speech_not_reach_fetch_enabled": true,
  "standardized_nudges_misinfo": true,
  "tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled": true,
  "rweb_video_timestamps_enabled": true,
  "longform_notetweets_rich_text_read_enabled": true,
  "longform_notetweets_inline_media_enabled": true,
  "responsive_web_enhance_cards_enabled": false,
  "spaces_2022_h2_clipping": false,
  "spaces_2022_h2_spaces_communities": false
}
```
If the request was successful, there is JSON response containing data about the Space. This JSON response contains a `media_key` under `data.audioSpace.metadata`. That `media_key`is important because it basically allows Twitter to access the M3U8 file.

If you're wondering what the heck an M3U8 file is, it is basically a plain text file that stores the URL paths of streaming audio or video and media track information. For example, one M3U8 file may give you references to online files for an YouTube video. I stole that definition from Google, you're welcome.

Okay, so now that the `media_key` is known, Twitter makes another GET request to `https://x.com/i/api/1.1/live_video_stream/status/<media_key>`, replacing `<media_key>` with the actual media key retrieved earlier. This returns another JSON response that looks similar to this:
```json
{
  "source" : {
    "location" : "https://example.com/audio/playlist.m3u8",
    "noRedirectPlaybackUrl" : "https://example.com/audio/playlist.m3u8",
    "status" : "LIVE_PUBLIC",
    "streamType" : "HLS"
  },
  "sessionId" : "1234567890123456789",
  "chatToken" : "dummyChatToken",
  "lifecycleToken" : "dummylifecycleToken",
  "shareUrl" : "https://twitter.com/i/broadcasts/12345"
}
```
The most interesting thing to take note of is `source.location`, because it holds the actual URL to the hosted M3U8 file. From here things become a bit easier. You can copy this URL and run the following:
```shell
yt-dlp <m3u8_url>
```

## Automating the entire thing with a CLI
What do you do when you can't access a Twitter Space without logging in and the official Twitter API costs nearly double your country's minimum wage? Naturally, you turn to web scraping. I started the project and created a `Downloader` class that will be instantiated later by the Command Line Interface (CLI).

```ts
export class Downloader implements DownloaderInterface {
  private username: string;
  private password: string;
  private options!: DownloaderOptions;
  private headers: TaskHeaders;
  private audioSpaceData!: Record<string, any>;
  private mediaKey!: string;
  private id: string;
  private isLoggedIn: boolean = false;
  private $: any;
  private playlist!: string;
  private playlistUrl!: string;
  private chunkBaseUrl!: string;
  private downloadChunksCount: number = 0;
  private storagePath;
  private chunksUrls!: string[];

  constructor(options: DownloaderOptions) {
    this.options = options;
    this.username = options.username;
    this.password = options.password;
    this.id = options.id;
    this.storagePath = path.resolve(`./task-${this.id}/`);

    this.headers = {
      'User-Agent': 'curl/7.81.0',
      'accept': "*/*",
      Referer: 'https://twitter.com/',
      'Content-Type': 'application/json',
    };

  }
  async init(): Promise<Downloader> {
    print.info("Starting authentication flow")
    await this.login();
    print.info(`Retrieving space metadata: [${this.id}]`);
    await this.setSpaceMetadataAndMediaKey();
    const playListInfoResponse: AxiosResponse = await getRequest(CONSTANTS.PLAYLIST_INFO_URL(this.mediaKey), this.headers);
    this.playlistUrl = playListInfoResponse.data.source.location;
    this.chunkBaseUrl = this.playlistUrl.replace(path.basename(this.playlistUrl), '');
    return this;
  }
  
  private async saveToDisk(data: any, location: string) {
    await fs.outputFile(path.join(this.storagePath + '/' + location), data);
  }
  
  async cleanup() {
    print.info("Cleaning up!");
    print.success("Done!");
  }
```

The `init` method is where the login happens, alongside other things like setting the media key and the M3U8 playlist URL. 

Twitter has this (weird?) login system, where logging in involves completing a series of subtasks. Each subtask returns a `flow_token` that must be provided when starting the next subtask.

To automate the login system, a guest token is required. The `getGuestToken` method scrapes the Twitter login page using Cheerio and retrieves the `<script>` element that sets the cookie. The guest token (`gt=`) is then extracted using regular expressions:

```ts
  private async getGuestToken(): Promise<string> {
    let scriptText = '';
    this.$('script').each((_: number, element: cheerio.Element) => {
      let text = this.$(element).html();
      if (text && text.includes('document.cookie')) {
        scriptText = text;
        return false; // Break the loop
      }
    });

    const stringWithGT: RegExpMatchArray | null = scriptText.match(/"gt=\d{19}/);
    if (stringWithGT && stringWithGT[0]) return stringWithGT[0].replace('"gt=', '');
    throw new Error('Failed to get guest token');
  }
```

`getGuestToken` is used in the `login` method like so. This code logins into the account by performing these subtasks in the sequence: `'' -> LoginJsInstrumentationSubtask -> LoginEnterUserIdentifierSSO -> AccountDuplicationCheck`. If the operation is successful, I will get an authentication token and a CSRF token [^1].

```ts
  async login() {
    try {
      const response: AxiosResponse = await getRequest(CONSTANTS.URL_BASE, this.headers);
      this.$ = cheerio.load(response.data);
      print.info("Retrieving guest token...");
      this.headers['X-Guest-Token'] = await this.getGuestToken();
      this.headers["Authorization"] = CONSTANTS.BEARER;
      // Initialize login flow:
      print.info('Logging in with credentials. Make sure 2FA is disabled on your account');
      let taskResponse: any;
      let taskInputs: any = CONSTANTS.LOGIN_FLOW_SUBTASK_DATA[''].input;
      taskResponse = (await postRequest(CONSTANTS.URL_FLOW_1, this.headers, JSON.stringify(taskInputs)));
      let flowToken: string = taskResponse.data.flow_token;
      // console.log(taskResponse.data)
      let nextSubtask: string = taskResponse.data.subtasks[0].subtask_id;

      const att: string = taskResponse.headers
        .get('set-cookie')
        .find((x: string) => x.startsWith('att='))
        .split('att=')[1]
        .split(';')[0];
      this.setHeaders({ cookie: `att=${att}` });


      while (!this.isLoggedIn) {
        if (Object.keys(CONSTANTS.LOGIN_FLOW_SUBTASK_DATA).find(x => x === nextSubtask)) {
          print.info(`Performing next subtask: ${nextSubtask}`);
        } else {
          throw new Error("Failed to get next subtask");
        }

        if (nextSubtask === 'LoginJsInstrumentationSubtask') {
          taskInputs = { flow_token: flowToken, ...CONSTANTS.LOGIN_FLOW_SUBTASK_DATA[nextSubtask].input };
          taskResponse = await postRequest(CONSTANTS.URL_FLOW_2, this.headers, JSON.stringify(taskInputs));
          flowToken = taskResponse.data.flow_token;
          nextSubtask = taskResponse.data.subtasks[0].subtask_id;
        } else if (nextSubtask === 'LoginEnterUserIdentifierSSO') {
          print.default('Submitting username...');
          taskInputs = { flow_token: flowToken, ...CONSTANTS.LOGIN_FLOW_SUBTASK_DATA[nextSubtask](this.username).input }
          taskResponse = await postRequest(CONSTANTS.URL_FLOW_2, this.headers, JSON.stringify(taskInputs));
          flowToken = taskResponse.data.flow_token;
          // console.log(taskResponse.data)
          nextSubtask = taskResponse.data.subtasks[0].subtask_id;
        } else if (nextSubtask === 'LoginEnterPassword') {
          print.default('Submitting password...');
          taskInputs = { flow_token: flowToken, ...CONSTANTS.LOGIN_FLOW_SUBTASK_DATA[nextSubtask](this.password).input }
          taskResponse = await postRequest(CONSTANTS.URL_FLOW_2, this.headers, JSON.stringify(taskInputs));
          flowToken = taskResponse.data.flow_token;
          nextSubtask = taskResponse.data.subtasks[0].subtask_id;
        } else if (nextSubtask === 'AccountDuplicationCheck') {
          print.info('Performing account duplication check')
          taskInputs = { flow_token: flowToken, ...CONSTANTS.LOGIN_FLOW_SUBTASK_DATA[nextSubtask].input };
          taskResponse = await postRequest(CONSTANTS.URL_FLOW_2, this.headers, JSON.stringify(taskInputs));
          flowToken = taskResponse.data.flow_token;
          nextSubtask = taskResponse.data.subtasks[0].subtask_id;
          this.isLoggedIn = true;
        }
      }

      print.info("Getting Authentication Token...");
      const twitterAuthToken = taskResponse.headers.get('set-cookie')
        .find((x: string) => x.startsWith('auth_token='))
        .split('auth_token=')[1]
        .split(';')[0];
      print.info("Getting CSRF Token...");

      let csrfToken = taskResponse.headers
        .get('set-cookie')
        .find((x: string) => x.startsWith('ct0='))
        .split('ct0=')[1]
        .split(';')[0];

      this.setHeaders({ cookie: `${this.headers.cookie}; auth_token=${twitterAuthToken}; ct0=${csrfToken}`, 'X-Csrf-Token': csrfToken });
      print.success("Login Success!\n\n");
      // console.log(this.headers);
    } catch (error) {
      throw error;
    }

    this.isLoggedIn = true;
  }
```

If nothing goes wrong, the `generateAudio` method is invoked. 
```ts
  async generateAudio() {
    this.playlist = await this.getPlaylist();
    this.chunksUrls = this.parsePlaylist();
    await this.downloadSegments(this.chunksUrls);
    await this.convertSegmentsToMp3();
    this.audioGenerated = true;
  }
```
In this method, I call `getPlaylist` to retrieve the playlist and store it to prevent unnecessary network calls [^2].
```ts
 private async getPlaylist() {
    let playlistPath: string = path.join(this.storagePath + "/" + "playlist.m3u8");
    let playlist: string;

    if (await fs.pathExists(playlistPath)) {
      print.info('Playlist already downloaded!');
      return await fs.readFile(playlistPath, { encoding: "utf-8" });
    }

    print.info('Downloading playlist');

    playlist = (await getRequest(this.playlistUrl, this.headers)).data;
    await this.saveToDisk(playlist, `playlist.m3u8`);
    return playlist;
  }
```

Next, the downloaded playlist is parsed with the [m3u8-parser package](https://www.npmjs.com/package/m3u8-parser). 
```ts
  private parsePlaylist(): string[] {
    const parser = new m3u8Parser.Parser();
    parser.push(this.playlist);
    parser.end();
    this.playlistManifest = parser.manifest;
    return parser.manifest.segments.map((x: { uri: string }) => this.chunkBaseUrl + x.uri);
  }
```
This makes it easier to retrieve the URL for each audio chunk. Next, each chunk URL is fed into `downloadSegments` where they are downloaded to the disk.
```ts
  private async downloadSegments(
    chunks: string[],
    retryCount: Record<string, number> = {},
    maxRetries: number = 10
  ): Promise<void> {

    // Check cache for the downloaded chunks

    print.info('Starting to download audio chunks')
    for (let url of chunks) {
      let message = `Starting to download chunks`;
      const chunkName = path.basename(url);
      const chunkStorageLocation: string = path.join('chunks', chunkName);
      if (!retryCount[chunkName]) retryCount[chunkName] = 0;
      if (await fs.pathExists(path.resolve(this.storagePath + "/" + chunkStorageLocation))) {
        this.downloadChunksCount++;
        message = `Skipping ${chunkName}`;
        // print.info(`${urlPath} already downloaded! Skipped!`);
        // break;
      } else {
        try {
          message = `Downloading ${chunkName}`
          const response = Buffer.from((await axios.get(url, { responseType: 'arraybuffer' })).data);
          this.downloadChunksCount++;
          // console.log(`Downloaded ${urlPath} ........................................ ${((this.downloadChunksCount / this.chunksUrls.length) * 100).toFixed(2)}% done`);
          await this.saveToDisk(response, chunkStorageLocation);
          // return response;
        } catch (error: any) {
          if (retryCount[chunkName] >= maxRetries) {
            throw new Error(`\nFailed to fetch chunk: ${chunkName}. Giving up after ${maxRetries} retries. \n${error.message}`);
          }

          retryCount[chunkName] += 1;
          console.error(`Failed to fetch ${chunkName} .................................. Retrying [${retryCount[chunkName]}/${maxRetries}]`);
          return this.downloadSegments([url], retryCount, maxRetries);
        }
      }
      print.progress(this.downloadChunksCount, this.chunksUrls.length, message, "AUDIO");
    }
  }
```

Finally, `convertSegmentsToMp3` converts each chunk to the `.mp3` audio format and combines them. `convertSegmentsToMp3` uses the [fluent-ffmpeg](https://www.npmjs.com/package/fluent-ffmpeg) Node.js FFmpeg wrapper to work with the audio files. This means that FFmpeg must be installed on your machine for it to work properly.

I read each downloaded chunk and wrote it into the PassThrough Duplex stream. With this, fluent-ffmpeg can read the stream change the format from `.aac` to `.mp3`, setting the audio codec as `libmp3lame` with a 44.1Khz sample rate frequency.

```ts
  private async convertSegmentsToMp3() {
    await fs.ensureDir(path.join(this.storagePath, 'out/'));
    const passThroughStream = new PassThrough();
    const finalOutputFilePath = path.join(this.storagePath, 'out/', `${this.audioSpaceData.metadata.title}.mp3`);
    // const ffmpegCommand = ffmpeg();
    const chunks: string[] = await fs.readdir(path.join(this.storagePath, 'chunks'), { encoding: "utf-8" });
    if (chunks.length === 0) {
      throw new Error('Failed to fetch chunks saved on disk.');
    }
    for (const chunkPath of chunks) {
      passThroughStream.write(await fs.readFile(path.join(this.storagePath, 'chunks/', chunkPath)));
      // ffmpegCommand.input(path.join(this.storagePath, 'chunks/', chunkPath))
    };
    passThroughStream.end();

    await new Promise<void>((resolve, reject) => {
      ffmpeg(passThroughStream)
        .inputFormat('aac')
        .audioFrequency(44100)  // Set sample rate to 44.1 kHz for better quality
        .audioChannels(2)       // Set audio channels to stereo
        .audioCodec('libmp3lame') // Set audio codec to libmp3lame for MP3 encoding
        .toFormat('mp3')        // Set output format to mp3
        .on('error', (err) => {
          reject(`Error: ${err.message}`);
        })
        .on('progress', (progress) => {
          const duration: number = new Date(Number(this.audioSpaceData.metadata.ended_at) - this.audioSpaceData.metadata.started_at).getTime();
          const datedTimeStamp: number = new Date(`1970-01-01T${progress.timemark}Z`).getTime();
          print.progress(datedTimeStamp, duration, "Combining chunks and converting to .mp3", "FFMPEG");
        })
        .on('end', () => {
          resolve();
          print.success('Merging completed');
        })
        .save(finalOutputFilePath);
    });
  }
```

With the `Downloader` class set up, all that's left is to create a CLI that will instantiate the class. This CLI uses the [Commander](https://www.npmjs.com/commander) package, although I feel it is a bit overkill for such a small project. But I am keeping it anyways.

```ts
import { Command } from 'commander';
import { Downloader } from '../index.js';
import { DownloaderOptions } from '../types.js';
const program = new Command();


program
    .name('spaces-dl')
    .description('CLI to download recorded Twitter Spaces')
    .version('0.0.1')
    .option('-m, --m3u8', 'm3u8 file url')
    .option('-i, --id <id>', 'A valid ID for a recorded Twitter Space')
    .option('-u, --username <username>', 'a valid twitter username without the @')
    .option('-p, --password <password>', 'a valid password for the username')
    .option('-o, --output <path>', 'output path for the recorded audio/video')
    .action((options) => {
        if (!options.id && !options.m3u8) {
            console.error("Error: --id option required");
            process.exit(1);
        }
    });

program.parse(process.argv);

const options: DownloaderOptions = program.opts();

try {
    console.log(options);
    let task: Downloader;
    task = await new Downloader(options).init();
    await task.generateAudio();
    await task.cleanup();
} catch (error) {
    throw error;
}
```

When the code is transpiled. In my terminal, I can run the following from the project root:
```shell
node dist/cli/cli.js --username <username> --password <password> -i <twitter_space_id>
```

The `twitter_space_id` can be obtained from the URL of the Space. For example, in https://x.com/i/spaces/1OyKAWgagpqJb, the ID is 1OyKAWgagpqJb.

![](https://i.imgur.com/AOx5txi.png)
Agba developer!

## Future plans
In the future, I hope to turn this into a full NPM package that you can install globally by running a command like `npm install -g spaces-dl`. I also want to implement alternative login systems in case the default login system with username and password fails. 

Since this is a side project, there are some 'missing' features, like the fact that you cannot provide a custom path for the output audio. I also want to improve the path system to work on Windows machines. You can find the full source code on [GitHub](https://github.com/daviduzondu/spaces-dl).

[^1]: One limitation of this approach is that sometimes, Twitter suspects that you are a bot, thus breaking the entire subtask flow system. So instead of the next subtask after `LoginJsInstrumentationSubtask` to be `LoginEnterUserIdentifierSSO`, I get `LoginEnterAlternateIdentifierSubtask` or sometimes `ArkoseLogin`. A potential workaround is providing the authentication token and CSRF token directly through the CLI, but I haven't gotten around to implementing that.

[^2]:  In cases where the user quits the program mid-operation and starts the operation later, this allows the program to use the stored M3U8 playlist instead of re-downloading it from the server. This storage mechanism also works for the audio chunks as well. 