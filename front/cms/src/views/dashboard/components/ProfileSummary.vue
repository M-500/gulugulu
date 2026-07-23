<template>
  <section class="profile-summary">
    <div class="profile-summary__identity">
      <div class="profile-summary__avatar">
        咕
        <span />
      </div>
      <div>
        <div class="profile-summary__name">
          <h2>咕噜创作者</h2>
          <span>已认证</span>
        </div>
        <p>@gulugulu · ID 100086</p>
        <small>分享生活灵感与实用创作技巧，让每一篇内容都有价值。</small>
      </div>
    </div>

    <div class="profile-summary__status">
      <div>
        <span>账号状态</span>
        <strong><i /> 正常</strong>
      </div>
      <div>
        <span>内容信用分</span>
        <strong>98 / 100</strong>
      </div>
      <button type="button">编辑资料</button>
    </div>
  </section>
</template>

<style scoped>
.profile-summary {
  display: flex;
  min-height: 128px;
  align-items: center;
  justify-content: space-between;
  gap: 28px;
  overflow: hidden;
  border: 1px solid #ebecef;
  border-radius: 14px;
  background:
    radial-gradient(circle at 92% 12%, rgba(239, 77, 88, 0.08), transparent 26%),
    #fff;
  padding: 22px 24px;
}

.profile-summary__identity {
  display: flex;
  min-width: 0;
  align-items: center;
  gap: 16px;
}

.profile-summary__avatar {
  position: relative;
  display: grid;
  width: 70px;
  height: 70px;
  flex: 0 0 70px;
  place-items: center;
  border: 4px solid #fff;
  border-radius: 22px;
  background: linear-gradient(145deg, #ff858b, #de3542);
  box-shadow: 0 8px 24px rgba(222, 53, 66, 0.2);
  color: #fff;
  font-size: 25px;
  font-weight: 800;
}

.profile-summary__avatar span {
  position: absolute;
  right: -1px;
  bottom: -1px;
  width: 15px;
  height: 15px;
  border: 3px solid #fff;
  border-radius: 50%;
  background: #12b76a;
}

.profile-summary__name {
  display: flex;
  align-items: center;
  gap: 9px;
}

.profile-summary__name h2 {
  margin: 0;
  color: #182230;
  font-size: 19px;
}

.profile-summary__name span {
  border-radius: 5px;
  background: #fff0f1;
  color: #e1434e;
  font-size: 10px;
  font-weight: 700;
  padding: 3px 6px;
}

.profile-summary__identity p {
  margin: 5px 0 6px;
  color: #7b8494;
  font-size: 12px;
}

.profile-summary__identity small {
  display: block;
  overflow: hidden;
  max-width: 480px;
  color: #98a2b3;
  font-size: 12px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.profile-summary__status {
  display: flex;
  flex: 0 0 auto;
  align-items: center;
  gap: 28px;
}

.profile-summary__status > div {
  display: grid;
  gap: 6px;
}

.profile-summary__status span {
  color: #98a2b3;
  font-size: 11px;
}

.profile-summary__status strong {
  color: #344054;
  font-size: 13px;
}

.profile-summary__status i {
  display: inline-block;
  width: 7px;
  height: 7px;
  margin-right: 4px;
  border-radius: 50%;
  background: #12b76a;
}

.profile-summary__status button {
  height: 36px;
  border: 1px solid #e4e7ec;
  border-radius: 8px;
  background: #fff;
  color: #475467;
  cursor: pointer;
  font-size: 12px;
  padding: 0 13px;
}

@media (max-width: 1040px) {
  .profile-summary {
    align-items: flex-start;
    flex-direction: column;
  }

  .profile-summary__status {
    width: 100%;
    justify-content: space-between;
    padding-left: 86px;
  }
}

@media (max-width: 620px) {
  .profile-summary {
    padding: 18px;
  }

  .profile-summary__avatar {
    width: 56px;
    height: 56px;
    flex-basis: 56px;
    border-radius: 18px;
  }

  .profile-summary__identity small {
    max-width: 220px;
  }

  .profile-summary__status {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 16px;
    padding-left: 0;
  }

  .profile-summary__status button {
    grid-column: 1 / -1;
  }
}
</style>
