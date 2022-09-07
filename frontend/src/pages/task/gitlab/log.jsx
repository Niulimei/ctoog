import './log.css';

export default ({ id }) => {
  const a = `[2022-09-07 11:17:08] [checkout] - git source enterprise__Tomato-Scrum/repo-web-page
  [2022-09-07 11:17:08] Stage [1/2] Starting to save ssh ... 
  [2022-09-07 11:17:08] Stage [2/2] Starting to download code ... 
  [2022-09-07 11:17:08] export LC_ALL=en_US.UTF-8 && rm -rf /root/workspace/enterprise__Tomato-Scrum/repo-web-page && mkdir -p /root/workspace/enterprise__Tomato-Scrum/repo-web-page && cd /root/workspace/enterprise__Tomato-Scrum/repo-web-page && git config --global user.name pipeline-temp && git config --global user.email pipeline-temp && git config --global http.postBuffer 1048576000 && git config --global http.lowSpeedLimit 0 && git config --global http.lowSpeedTime 999999 && git config --global core.compression -1 && git config --global ssh.postBuffer 1024M && git config --global ssh.maxRequsetBuffer 1024M && git clone --depth 5 ssh://git@osc.gitee.work:5503/enterprise__Tomato-Scrum/repo-web-page.git -b bugfix_pulledit . && git reset --hard cc90b30124afe77f2d90f41936de640d4845a5db || (git pull --unshallow && git reset --hard cc90b30124afe77f2d90f41936de640d4845a5db)
  [2022-09-07 11:17:08] Cloning into '.'...
  [2022-09-07 11:17:08] 
  [2022-09-07 11:17:09] Warning: Permanently added '[osc.gitee.work]:5503,[172.20.0.250]:5503' (ECDSA) to the list of known hosts.
  [2022-09-07 11:17:09] 
  [2022-09-07 11:17:11] HEAD is now at cc90b30 code-plugins-2346 编辑PR接口改为go
  [2022-09-07 11:17:11] 
  [2022-09-07 11:17:11] All Stage Completed Successfully
  [2022-09-07 11:17:11] The Repo Clone Successfully！`
  return (
    <pre className={"terminal"}>
      <code>
      {a}
      </code>
    </pre>
  );
}