Today my first ever merge request to a free and open source project was
merged :O). You can check it out here:
[(!142) - alpine / mkinitfs - GitLab](https://gitlab.alpinelinux.org/alpine/mkinitfs/-/merge_requests/142?commit_id=33830906afb5641dfe97162e1d348935fc14f641).
So what did I change and why did I change it?

A while back I tried to set up hibernation support for my laptop
running alpine linux. For some reason I did not have and did not
want a swap partition. But if the whole state of the system is
supposed to be saved to disk, it has to go somewhere. So I created a
swap file and tried to hibernate to that. But it did not work.

I started digging. What was supposed to happen? After telling the kernel
to hibernate, it stops the system, saves the state of the system to swap
(in my case the swap file) and powers off the computer. Then after
turning the computer back on, the bootloader has to specify the
location of the system state on the kernel command line. Either as
reference to a swap partition or as a reference to a partition with a
regular file system and the offset of the swap file on that partition.
These parameters are then handled in the early userspace. There are two
great blog posts that explain how the early userspace works:

 - [Messing with your initramfs - Hare edition](https://bitfehler.srht.site/posts/2022-09-29_messing-with-your-initramfs---hare-style-.html)
 - [Messing with your initramfs - Alpine edition](https://bitfehler.srht.site/posts/2022-11-28_messing-with-your-initramfs---alpine-edition.html)

After looking through alpines initramfs-init script I realized, that
alpine just ignores the `resume_offset` parameter. So I changed that.
I only had to add three lines of code and now its working. In order to
get my changes merged I also had to document them on the man page and
write a test case. It turns out writing the test case was the hardest
part, because the test framework confused me.

Either way I am really excited that I contributed to a free and open
source project for the first time.
