# Copyright (c) seasonjs. All rights reserved.
# Licensed under the MIT License. See License.txt in the project root for license information.


from huggingface_hub import hf_hub_download


def main():
    hf_hub_download(
        repo_id="julien-c/dummy-unknown",
        filename="config.json",
        revision="main",
        cache_dir="../tmp",
        resume_download=True,
    )


if __name__ == "__main__":
    main()
